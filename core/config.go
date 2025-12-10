package core

import (
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
)

type AppConfig struct {
	DbPath             string
	BaseUrl            string
	JWTSecret          []byte
	JWTExpiration      time.Duration
	JWTCookieAllowHTTP bool
	AppBuildVersion    string
	AppBuildDate       string
	AppBuildCommit     string
	AppGinMode         string
	AppPort            string
	AppUsersToCreate   []User
	AppUserPattern     *regexp.Regexp
	AppKeyPattern      *regexp.Regexp
	AppDataMaxSize     int64
	AppKeysPerUser     int64
	SwaggerEnabled     bool
}

var Config = func() AppConfig {
	config := AppConfig{
		DbPath:             resolvePath(os.Getenv("GENESIS_DB_PATH")),
		BaseUrl:            os.Getenv("GENESIS_BASE_URL"),
		JWTSecret:          []byte(os.Getenv("GENESIS_JWT_SECRET")),
		JWTExpiration:      time.Duration(parseInt(os.Getenv("GENESIS_JWT_TOKEN_EXPIRATION"))) * time.Minute,
		JWTCookieAllowHTTP: os.Getenv("GENESIS_JWT_COOKIE_ALLOW_HTTP") == "true",
		AppBuildVersion:    os.Getenv("GENESIS_BUILD_VERSION"),
		AppBuildDate:       os.Getenv("GENESIS_BUILD_DATE"),
		AppBuildCommit:     os.Getenv("GENESIS_BUILD_COMMIT"),
		AppGinMode:         os.Getenv("GENESIS_GIN_MODE"),
		AppPort:            os.Getenv("GENESIS_PORT"),
		AppUsersToCreate:   parseInitialUserList(os.Getenv("GENESIS_CREATE_USERS")),
		AppUserPattern:     regexp.MustCompile(os.Getenv("GENESIS_USERNAME_PATTERN")),
		AppKeyPattern:      regexp.MustCompile(os.Getenv("GENESIS_KEY_PATTERN")),
		AppDataMaxSize:     parseInt(os.Getenv("GENESIS_DATA_MAX_SIZE")) * 1000,
		AppKeysPerUser:     parseInt(os.Getenv("GENESIS_KEYS_PER_USER")),
		SwaggerEnabled:     os.Getenv("GENESIS_SWAGGER_ENABLED") != "false", // Enabled by default
	}

	Logger.Debug("build info",
		zap.String("version", config.AppBuildVersion),
		zap.String("date", config.AppBuildDate),
		zap.String("commit", config.AppBuildCommit),
	)

	return config
}()

func parseInitialUserList(raw string) []User {
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
				Name:     strings.TrimSuffix(user[0], "!"),
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
	workdir, _ := os.Getwd()
	return workdir
}
