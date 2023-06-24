package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgraph-io/badger/v4"
	"golang.org/x/crypto/bcrypt"
	"log"
)

const (
	DbUserPrefix     = "user:"
	DbUserDataPrefix = "data:"
)

var database *badger.DB

type User struct {
	User     string `json:"user"`     // Username
	Password string `json:"password"` // Hashed password
}

func CreateUser(name string, password string) error {
	txn := database.NewTransaction(true)
	key := []byte(DbUserPrefix + name)

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

func AuthenticateUser(name string, password string) (*User, error) {
	user, err := GetUser(name)

	if err != nil {
		return nil, fmt.Errorf("failed to parse data: %v", err)
	} else if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return nil, errors.New("invalid password")
	}

	return user, nil
}

func GetUser(name string) (*User, error) {
	txn := database.NewTransaction(false)
	key := []byte(DbUserPrefix + name)

	data, err := txn.Get(key)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve data: %v", err)
	}

	var user User
	return &user, data.Value(func(val []byte) error {
		return json.Unmarshal(val, &user)
	})
}

func SetDataForUser(name string, key string, data map[string]interface{}) error {
	txn := database.NewTransaction(true)

	if data, err := json.Marshal(data); err != nil {
		return err
	} else if err := txn.Set([]byte(DbUserDataPrefix+name+":"+key), data); err != nil {
		return err
	} else {
		return txn.Commit()
	}
}

func DeleteDataFromUser(name string, key string) error {
	txn := database.NewTransaction(true)

	if err := txn.Delete([]byte(DbUserDataPrefix + name + ":" + key)); err != nil {
		return err
	} else {
		return txn.Commit()
	}
}

func GetDataFromUser(name string, key string) ([]byte, error) {
	txn := database.NewTransaction(false)
	item, err := txn.Get([]byte(DbUserDataPrefix + name + ":" + key))

	if err != nil {
		return nil, err
	}

	var data []byte
	return data, item.Value(func(v []byte) error {
		*&data = v
		return nil
	})
}

func GetAllDataFromUser(name string) ([]byte, error) {
	txn := database.NewTransaction(false)
	it := txn.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()

	prefix := []byte(DbUserDataPrefix + name + ":")
	data := make(map[string]interface{}, 0)

	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		item := it.Item()
		key := item.Key()

		err := item.Value(func(v []byte) error {
			var raw map[string]interface{}

			if err := json.Unmarshal(v, &raw); err != nil {
				return err
			} else {
				data[string(key[len(prefix):])] = raw
			}

			return nil
		})

		if err != nil {
			break
		}
	}

	return json.Marshal(data)
}

func DropDatabase() {
	if database.DropAll() != nil {
		log.Fatal("Failed to drop database")
	}
}

func init() {
	options := badger.DefaultOptions(Config().DbPath)

	if db, err := badger.Open(options); err != nil {
		log.Fatal(err)
	} else {
		database = db
	}
}
