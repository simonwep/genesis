package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgraph-io/badger/v4"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

const (
	DbUserPrefix         = "user:"
	DbUserDataPrefix     = "data:"
	DbExpiredTokenPrefix = "expired:"
)

var database *badger.DB

type User struct {
	User     string `json:"user"`     // Username
	Password string `json:"password"` // Hashed password
}

func CreateUser(name string, password string) error {
	txn := database.NewTransaction(true)
	defer txn.Discard()

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
		return nil, err
	} else if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return nil, errors.New("invalid password")
	}

	return user, nil
}

func GetUser(name string) (*User, error) {
	txn := database.NewTransaction(false)
	defer txn.Discard()

	key := []byte(DbUserPrefix + name)

	data, err := txn.Get(key)
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return nil, nil
		} else {
			return nil, fmt.Errorf("failed to retrieve data: %v", err)
		}
	}

	var user User
	return &user, data.Value(func(val []byte) error {
		return json.Unmarshal(val, &user)
	})
}

func SetPasswordForUser(name string, password string) error {
	user, err := GetUser(name)

	if err != nil {
		return err
	} else if user == nil {
		return fmt.Errorf("no such user with name %v", name)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	txn := database.NewTransaction(true)
	defer txn.Discard()

	key := []byte(DbUserPrefix + name)

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

func SetDataForUser(name string, key string, data []byte) error {
	txn := database.NewTransaction(true)
	defer txn.Discard()

	if err := txn.Set([]byte(DbUserDataPrefix+name+":"+key), data); err != nil {
		return err
	} else {
		return txn.Commit()
	}
}

func DeleteDataFromUser(name string, key string) error {
	txn := database.NewTransaction(true)
	defer txn.Discard()

	if err := txn.Delete([]byte(DbUserDataPrefix + name + ":" + key)); err != nil {
		return err
	} else {
		return txn.Commit()
	}
}

func GetDataFromUser(name string, key string) ([]byte, error) {
	txn := database.NewTransaction(false)
	defer txn.Discard()

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
	defer txn.Discard()

	it := txn.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()

	prefix := []byte(DbUserDataPrefix + name + ":")
	data := make([]string, 0)

	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		item := it.Item()
		key := item.Key()

		err := item.Value(func(v []byte) error {
			if rawKey, err := json.Marshal(string(key[len(prefix):])); err != nil {
				return err
			} else {
				data = append(data, string(rawKey)+":"+string(v))
			}

			return nil
		})

		if err != nil {
			break
		}
	}

	return []byte("{" + strings.Join(data, ",") + "}"), nil
}

func GetDataCountForUser(name, includedKey string) int64 {
	txn := database.NewTransaction(false)
	defer txn.Discard()

	it := txn.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()

	prefix := []byte(DbUserDataPrefix + name + ":")
	hadIncludedKey := false
	count := int64(0)

	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		if !hadIncludedKey {
			key := string(it.Item().Key())
			hadIncludedKey = key == DbUserDataPrefix+name+":"+includedKey
		}

		count++
	}

	if !hadIncludedKey {
		count++
	}

	return count
}

func StoreInvalidatedToken(jti string, expiration time.Duration) error {
	return database.Update(func(txn *badger.Txn) error {
		key := []byte(DbExpiredTokenPrefix + ":" + jti)
		return txn.SetEntry(badger.NewEntry(key, []byte{}).WithTTL(expiration))
	})
}

func IsTokenBlacklisted(jti string) (bool, error) {
	txn := database.NewTransaction(false)
	defer txn.Discard()

	item, err := txn.Get([]byte(DbExpiredTokenPrefix + ":" + jti))

	if err == badger.ErrKeyNotFound {
		return false, nil
	} else {
		return item != nil, err
	}
}

func ResetDatabase() {
	if err := database.DropAll(); err != nil {
		Fatal("failed to drop database", zap.Error(err))
	}

	initializeUsers()
}

func initializeUsers() {
	for _, user := range Config.AppUsersToCreate {
		usr, err := GetUser(user.Name)

		if err != nil {
			Error("failed to check for user", zap.Error(err))
		} else if usr == nil {
			if err = CreateUser(user.Name, user.Password); err != nil {
				Error("failed to create user", zap.Error(err))
			} else {
				Debug("created new user", zap.String("name", user.Name))
			}
		}
	}
}

func printDebugInformation() {
	txn := database.NewTransaction(false)
	defer txn.Discard()

	it := txn.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()

	results := make(map[string]int)
	results[DbUserPrefix] = 0
	results[DbUserDataPrefix] = 0
	results[DbExpiredTokenPrefix] = 0

	for it.Rewind(); it.Valid(); it.Next() {
		key := strings.Split(string(it.Item().Key()), ":")
		results[key[0]+":"]++
	}

	Debug("users", zap.Int("count", results[DbUserPrefix]))
	Debug("datasets", zap.Int("count", results[DbUserDataPrefix]))
	Debug("expired keys", zap.Int("count", results[DbExpiredTokenPrefix]))
}

func init() {
	options := badger.DefaultOptions(Config.DbPath)
	options.Logger = nil

	if db, err := badger.Open(options); err != nil {
		Fatal("failed to open database", zap.Error(err))
	} else {
		database = db
	}

	printDebugInformation()
	initializeUsers()
}
