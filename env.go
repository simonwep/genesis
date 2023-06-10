package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	GinMode      string
	DbPath       string
	AllowedUsers []string
	JWTSecret    []byte
	JWTExpires   time.Duration
}

var (
	env       Config
	envLoaded = false
)

func Env() *Config {
	if !envLoaded {
		loadEnvVariables()
		envLoaded = true
	}

	return &env
}

func loadEnvVariables() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	env = Config{
		GinMode:      os.Getenv("GIN_MODE"),
		DbPath:       os.Getenv("DB_PATH"),
		AllowedUsers: strings.Split(os.Getenv("ALLOWED_USERS"), ","),
		JWTSecret:    []byte(os.Getenv("JWT_SECRET")),
		JWTExpires:   time.Duration(parseInt(os.Getenv("JWT_EXPIRES_IN"))) * time.Minute,
	}
}

func parseInt(str string) int64 {
	if value, err := strconv.ParseInt(str, 10, 64); err != nil {
		panic(err)
	} else {
		return value
	}
}
