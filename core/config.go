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
	GinMode              string
	DbPath               string
	JWTSecret            []byte
	JWTExpires           time.Duration
	AppPort              string
	AppAllowedUsers      []string
	AppAllowedKeyPattern *regexp.Regexp
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
		AppAllowedUsers:      strings.Split(os.Getenv("APP_ALLOWED_USERS"), ","),
		AppAllowedKeyPattern: regexp.MustCompile(os.Getenv("APP_KEY_PATTERN")),
	}
}

func parseInt(str string) int64 {
	if value, err := strconv.ParseInt(str, 10, 64); err != nil {
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
