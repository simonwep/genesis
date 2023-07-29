package core

import (
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type AppUser struct {
	Name     string
	Password string
}

type AppConfig struct {
	DbPath               string
	JWTSecret            []byte
	JWTAccessExpiration  time.Duration
	JWTRefreshExpiration time.Duration
	AppGinMode           string
	AppLogMode           string
	AppPort              string
	AppInitialUsers      []AppUser
	AppAllowedKeyPattern *regexp.Regexp
	AppValueMaxSize      int64
	AppKeysPerUser       int64
}

var Config = func() AppConfig {
	if err := godotenv.Load(path.Join(currentDir(), ".env")); err != nil {
		log.Fatal("failed to retrieve data", zap.Error(err))
	}

	return AppConfig{
		DbPath:               resolvePath(os.Getenv("DB_PATH")),
		JWTSecret:            []byte(os.Getenv("JWT_SECRET")),
		JWTAccessExpiration:  time.Duration(parseInt(os.Getenv("JWT_ACCESS_TOKEN_EXPIRATION"))) * time.Minute,
		JWTRefreshExpiration: time.Duration(parseInt(os.Getenv("JWT_REFRESH_TOKEN_EXPIRATION"))) * time.Minute,
		AppGinMode:           os.Getenv("APP_GIN_MODE"),
		AppLogMode:           os.Getenv("APP_LOG_MODE"),
		AppPort:              os.Getenv("APP_PORT"),
		AppInitialUsers:      parseUserPasswordList(os.Getenv("APP_INITIAL_USERS")),
		AppAllowedKeyPattern: regexp.MustCompile(os.Getenv("APP_KEY_PATTERN")),
		AppValueMaxSize:      parseInt(os.Getenv("APP_VALUE_MAX_SIZE")) * 1000,
		AppKeysPerUser:       parseInt(os.Getenv("APP_KEYS_PER_USER")),
	}
}()

func parseUserPasswordList(raw string) []AppUser {
	list := make([]AppUser, 0)

	if len(raw) == 0 {
		return list
	}

	for _, item := range strings.Split(raw, ",") {
		user := strings.Split(item, ":")

		if len(user) != 2 {
			log.Fatal("invalid pattern for allowed users")
		}

		list = append(list, AppUser{
			Name:     user[0],
			Password: user[1],
		})
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
