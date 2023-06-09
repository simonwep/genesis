package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgraph-io/badger/v4"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/slices"
	"log"
	"os"
	"path"
	"strings"
)

var database *badger.DB

type User struct {
	User     string `json:"user"`     // Username
	Password string `json:"password"` // Hashed password
	Data     string `json:"data"`     // User data
}

func ValidateUserName(name string) bool {
	allowedUser := GetEnv("ALLOWED_USERS")
	return allowedUser == "" || slices.Contains(strings.Split(allowedUser, ","), name)
}

func CreateUser(name string, password string) error {
	if !ValidateUserName(name) {
		return fmt.Errorf("a user with the name %v cannot be created", name)
	}

	txn := database.NewTransaction(true)
	key := []byte(name)

	if item, err := txn.Get(key); item != nil {
		return fmt.Errorf("a user with the name %v already exists", name)
	} else if err != nil && err != badger.ErrKeyNotFound {
		return fmt.Errorf("failed to check if user already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	data, err := json.Marshal(User{
		User:     name,
		Password: string(hash),
		Data:     "{}",
	})

	if err != nil {
		return fmt.Errorf("failed to create user data: %v", err)
	} else if err := txn.Set(key, data); err != nil {
		return fmt.Errorf("failed to store user: %v", err)
	} else if err := txn.Commit(); err != nil {
		return fmt.Errorf("failed to commit data: %v", err)
	}

	return nil
}

func GetUser(name string, password string) (*User, error) {
	txn := database.NewTransaction(false)
	key := []byte(name)

	data, err := txn.Get(key)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve data: %v", err)
	}

	var user User
	err := data.Value(func(val []byte) error {
		return json.Unmarshal(val, &user)
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse data: %v", err)
	} else if !checkPassword(&user, password) {
		return nil, errors.New("invalid password")
	}

	return &user, nil
}

func checkPassword(user *User, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil
}

func init() {
	wd, err := os.Getwd()

	if err != nil {
		log.Fatal(err)
	}

	dir := path.Join(wd, GetEnv("DB_PATH"))
	options := badger.DefaultOptions(dir)

	if db, err := badger.Open(options); err != nil {
		log.Fatal(err)
	} else {
		database = db
	}
}
