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

type AppUser struct {
	Name     string
	Password string
}

type AppConfig struct {
	GinMode              string
	DbPath               string
	JWTSecret            []byte
	JWTExpires           time.Duration
	AppPort              string
	AppInitialUsers      []AppUser
	AppAllowedKeyPattern *regexp.Regexp
	AppValueMaxSize      int64
	AppKeysPerUser       int64
}

var Config AppConfig

func init() {
	if err := godotenv.Load(path.Join(currentDir(), ".env")); err != nil {
		Logger.Fatal("failed to retrieve data", zap.Error(err))
	}

	Config = AppConfig{
		GinMode:              os.Getenv("GIN_MODE"),
		DbPath:               resolvePath(os.Getenv("DB_PATH")),
		JWTSecret:            []byte(os.Getenv("JWT_SECRET")),
		JWTExpires:           time.Duration(parseInt(os.Getenv("JWT_EXPIRES_IN"))) * time.Minute,
		AppPort:              os.Getenv("APP_PORT"),
		AppInitialUsers:      parseUserPasswordList(os.Getenv("APP_INITIAL_USERS")),
		AppAllowedKeyPattern: regexp.MustCompile(os.Getenv("APP_KEY_PATTERN")),
		AppValueMaxSize:      parseInt(os.Getenv("APP_VALUE_MAX_SIZE")) * 1000,
		AppKeysPerUser:       parseInt(os.Getenv("APP_KEYS_PER_USER")),
	}
}

func parseUserPasswordList(raw string) []AppUser {
	list := make([]AppUser, 0)

	if len(raw) == 0 {
		return list
	}

	for _, item := range strings.Split(raw, ",") {
		user := strings.Split(item, ":")

		if len(user) != 2 {
			Logger.Fatal("invalid pattern for allowed users")
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
