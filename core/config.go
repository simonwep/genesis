package core

import (
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type AppConfig struct {
	DbPath               string
	JWTSecret            []byte
	JWTAccessExpiration  time.Duration
	JWTRefreshExpiration time.Duration
	AppGinMode           string
	AppLogMode           string
	AppPort              string
	AppUsersToCreate     []User
	AppKeyPattern        *regexp.Regexp
	AppDataMaxSize       int64
	AppKeysPerUser       int64
}

var Config = func() AppConfig {
	if err := godotenv.Load(path.Join(currentDir(), ".env")); err != nil {
		Logger.Warn("failed to retrieve data", zap.Error(err))
	}

	return AppConfig{
		DbPath:               resolvePath(os.Getenv("GENESIS_DB_PATH")),
		JWTSecret:            []byte(os.Getenv("GENESIS_JWT_SECRET")),
		JWTAccessExpiration:  time.Duration(parseInt(os.Getenv("GENESIS_JWT_ACCESS_TOKEN_EXPIRATION"))) * time.Minute,
		JWTRefreshExpiration: time.Duration(parseInt(os.Getenv("GENESIS_JWT_REFRESH_TOKEN_EXPIRATION"))) * time.Minute,
		AppGinMode:           os.Getenv("GENESIS_GIN_MODE"),
		AppLogMode:           os.Getenv("GENESIS_LOG_MODE"),
		AppPort:              os.Getenv("GENESIS_PORT"),
		AppUsersToCreate:     parseUserPasswordList(os.Getenv("GENESIS_CREATE_USERS")),
		AppKeyPattern:        regexp.MustCompile(os.Getenv("GENESIS_KEY_PATTERN")),
		AppDataMaxSize:       parseInt(os.Getenv("GENESIS_DATA_MAX_SIZE")) * 1000,
		AppKeysPerUser:       parseInt(os.Getenv("GENESIS_KEYS_PER_USER")),
	}
}()

func parseUserPasswordList(raw string) []User {
	list := make([]User, 0)

	if len(raw) == 0 {
		return list
	}

	for _, item := range strings.Split(raw, ",") {
		user := strings.Split(item, ":")

		if len(user) != 2 {
			Logger.Warn("invalid pattern for allowed users", zap.String("user", item))
		} else {
			list = append(list, User{
				User:     strings.TrimSuffix(user[0], "!"),
				Admin:    strings.HasSuffix(user[0], "!"),
				Password: user[1],
			})
		}
	}

	return list
}

func parseInt(str string) int64 {
	raw := strings.ReplaceAll(str, "_", "")
	if value, err := strconv.ParseInt(raw, 10, 64); err != nil {
		panic(err)
	} else {
		return value
	}
}

func resolvePath(path string) string {
	return filepath.Join(currentDir(), path)
}

func currentDir() string {
	_, filename, _, _ := runtime.Caller(0)
	return path.Join(path.Dir(filename), "..")
}
