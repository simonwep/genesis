package core

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgraph-io/badger/v4"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

const (
	dbKeySeparator       = "/"
	dbUserPrefix         = "usr" // user:{name}
	dbDataPrefix         = "dat"
	dbExpiredTokenPrefix = "exp" // data:{name}:{key}
)

var (
	ErrUserAlreadyExists = errors.New("a user with this name already exists")
	ErrUserNotFound      = errors.New("user not found")
)

// User represents a user in the system
// @Description User with credentials
type User struct {
	Name     string `json:"name" validate:"required,gte=3,lte=32" example:"admin"`
	Admin    bool   `json:"admin" example:"true"`
	Password string `json:"password" validate:"required,gte=8,lte=64" example:"password123"`
}

// PartialUser represents partial user data for updates
// @Description Partial user data (both fields optional)
type PartialUser struct {
	Admin    *bool   `json:"admin,omitempty" example:"false"`
	Password *string `json:"password,omitempty" validate:"omitempty,gte=8,lte=64" example:"newPassword123"`
}

// PublicUser represents user information without sensitive data
// @Description User information returned to clients (no password)
type PublicUser struct {
	Name  string `json:"name" example:"admin"`
	Admin bool   `json:"admin" example:"true"`
}

var database *badger.DB

func CreateUser(user User) error {
	txn := database.NewTransaction(true)
	key := buildUserKey(user.Name)
	defer txn.Discard()

	if existingUser, err := GetUser(user.Name); existingUser != nil {
		return ErrUserAlreadyExists
	} else if err != nil {
		return fmt.Errorf("failed to check if user already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	} else if data, err := json.Marshal(User{
		Name:     user.Name,
		Admin:    user.Admin,
		Password: string(hash),
	}); err != nil {
		return fmt.Errorf("failed to create user data: %w", err)
	} else if err := txn.Set(key, data); err != nil {
		return fmt.Errorf("failed to store user: %w", err)
	} else if err := txn.Commit(); err != nil {
		return fmt.Errorf("failed to commit data: %w", err)
	}

	return nil
}

func UpdateUser(name string, user PartialUser) error {
	txn := database.NewTransaction(true)
	key := buildUserKey(name)
	defer txn.Discard()

	existingUser, err := GetUser(name)
	if err != nil {
		if errors.Is(err, badger.ErrKeyNotFound) {
			return ErrUserNotFound
		}

		return fmt.Errorf("failed to check if user exists")
	}

	if user.Password == nil {
		user.Password = &existingUser.Password
	} else {
		if hash, err := hashPassword(*user.Password); err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		} else {
			user.Password = &hash
		}
	}

	if user.Admin == nil {
		user.Admin = &existingUser.Admin
	}

	if data, err := json.Marshal(User{
		Name:     name,
		Admin:    *user.Admin,
		Password: *user.Password,
	}); err != nil {
		return fmt.Errorf("failed to create user data: %w", err)
	} else if err := txn.Set(key, data); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	} else if err := txn.Commit(); err != nil {
		return fmt.Errorf("failed to commit data: %w", err)
	}

	return nil
}

func AuthenticateUser(name string, password string) (*User, error) {
	user, err := GetUser(name)

	if err != nil {
		return nil, err
	} else if user == nil {
		return nil, nil
	} else if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return nil, errors.New("invalid password")
	}

	return user, nil
}

func GetUser(name string) (*User, error) {
	txn := database.NewTransaction(false)
	key := buildUserKey(name)
	defer txn.Discard()

	data, err := txn.Get(key)
	if err != nil {
		if errors.Is(err, badger.ErrKeyNotFound) {
			return nil, nil
		} else {
			return nil, fmt.Errorf("failed to retrieve data: %w", err)
		}
	}

	var user User
	return &user, data.Value(func(val []byte) error {
		return json.Unmarshal(val, &user)
	})
}

func GetUsers(skip string) ([]*PublicUser, error) {
	txn := database.NewTransaction(true)
	defer txn.Discard()

	it := txn.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()

	skipKey := buildUserKey(skip)
	users := make([]*PublicUser, 0)
	prefix := buildUserKey("")

	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		item := it.Item()

		// Skip the user we want to skip
		if bytes.Equal(skipKey, item.Key()) {
			continue
		}

		var user PublicUser
		err := item.Value(func(val []byte) error {
			return json.Unmarshal(val, &user)
		})

		if err != nil {
			return nil, err
		}

		users = append(users, &user)
	}

	return users, nil
}

func GetAllUsers() ([]*PublicUser, error) {
	txn := database.NewTransaction(true)
	defer txn.Discard()

	it := txn.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()

	users := make([]*PublicUser, 0)
	prefix := buildUserKey("")

	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		item := it.Item()

		var user PublicUser
		err := item.Value(func(val []byte) error {
			return json.Unmarshal(val, &user)
		})

		if err != nil {
			return nil, err
		}

		users = append(users, &user)
	}

	return users, nil
}

func DeleteUser(name string) error {
	txn := database.NewTransaction(true)
	defer txn.Discard()

	it := txn.NewIterator(badger.DefaultIteratorOptions)

	// Remove data
	prefix := buildUserDataKey(name, "")
	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		if err := txn.Delete(it.Item().Key()); err != nil {
			it.Close()
			return err
		}
	}

	it.Close()

	// Remove user
	if err := txn.Delete(buildUserKey(name)); err != nil {
		return err
	}

	return txn.Commit()
}

func SetDataForUser(name string, key string, data []byte) error {
	txn := database.NewTransaction(true)
	defer txn.Discard()

	if err := txn.Set(buildUserDataKey(name, key), data); err != nil {
		return err
	} else {
		return txn.Commit()
	}
}

func DeleteDataFromUser(name string, key string) error {
	txn := database.NewTransaction(true)
	defer txn.Discard()

	if err := txn.Delete(buildUserDataKey(name, key)); err != nil {
		return err
	} else {
		return txn.Commit()
	}
}

func GetDataFromUser(name string, key string) ([]byte, error) {
	txn := database.NewTransaction(false)
	defer txn.Discard()

	item, err := txn.Get(buildUserDataKey(name, key))
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

	prefix := buildUserDataKey(name, "")
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

	prefix := buildUserDataKey(name, "")
	hadIncludedKey := false
	count := int64(0)

	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		if !hadIncludedKey {
			key := string(it.Item().Key())
			hadIncludedKey = key == string(buildUserDataKey(name, includedKey))
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
		return txn.SetEntry(badger.NewEntry(buildExpiredKey(jti), []byte{}).WithTTL(expiration))
	})
}

func IsTokenBlacklisted(jti string) (bool, error) {
	txn := database.NewTransaction(false)
	defer txn.Discard()

	item, err := txn.Get(buildExpiredKey(jti))

	if errors.Is(err, badger.ErrKeyNotFound) {
		return false, nil
	} else {
		return item != nil, err
	}
}

func ResetDatabase() {
	if err := database.DropAll(); err != nil {
		Logger.Fatal("failed to drop database", zap.Error(err))
	}

	InitializeUsers()
}

func InitializeUsers() {
	for _, user := range Config.AppUsersToCreate {
		if existingUser, err := GetUser(user.Name); err != nil {
			Logger.Error("failed to check for user", zap.Error(err))
		} else if existingUser != nil {
			continue
		}

		if err := CreateUser(user); err != nil {
			Logger.Error("failed to create user", zap.Error(err))
		} else {
			Logger.Info("created new user", zap.String("name", user.Name), zap.Bool("admin", user.Admin))
		}
	}
}

func printDebugInformation() {
	txn := database.NewTransaction(false)
	defer txn.Discard()

	it := txn.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()

	results := make(map[string]int)
	results[dbUserPrefix] = 0
	results[dbDataPrefix] = 0
	results[dbExpiredTokenPrefix] = 0

	for it.Rewind(); it.Valid(); it.Next() {
		key := strings.Split(string(it.Item().Key()), dbKeySeparator)
		results[key[0]]++
	}

	Logger.Debug("users", zap.Int("count", results[dbUserPrefix]))
	Logger.Debug("datasets", zap.Int("count", results[dbDataPrefix]))
	Logger.Debug("expired keys", zap.Int("count", results[dbExpiredTokenPrefix]))
}

func buildExpiredKey(key string) []byte {
	return []byte(dbExpiredTokenPrefix + dbKeySeparator + key)
}

func buildUserKey(name string) []byte {
	return []byte(dbUserPrefix + dbKeySeparator + name)
}

func buildUserDataKey(name, key string) []byte {
	return []byte(dbDataPrefix + dbKeySeparator + name + dbKeySeparator + key)
}

func hashPassword(pwd string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	} else {
		return string(hashed), err
	}
}

func init() {
	options := badger.DefaultOptions(Config.DbPath)
	options.Logger = nil

	// Adjust options for a smaller database
	options.CompactL0OnClose = true
	options.ValueLogFileSize = 64 << 20 // 64MB
	options.NumLevelZeroTables = 1
	options.NumLevelZeroTablesStall = 2

	if db, err := badger.Open(options); err != nil {
		Logger.Fatal("failed to open database", zap.Error(err))
	} else {
		database = db
	}

	// Shutdown database gracefully
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	go func() {
		sig := <-sigs
		Logger.Info("received signal, closing database", zap.String("signal", sig.String()))

		if err := database.Close(); err != nil {
			Logger.Error("failed to close database", zap.Error(err))
		}

		os.Exit(0)
	}()

	// Run garbage collector once an hour
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()

		for {
			<-ticker.C
			err := database.RunValueLogGC(0.5)
			if errors.Is(err, badger.ErrNoRewrite) {
				continue
			} else if err != nil {
				Logger.Error("failed to run value log GC", zap.Error(err))
			}
		}
	}()

	printDebugInformation()
}
