package core

import (
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"log"
	"os"
	"path"
	"testing"
)

var Logger = func() *zap.Logger {
	envFile := path.Join(currentDir(), ".env")
	if testing.Testing() {
		envFile = path.Join(currentDir(), ".env.test")
	}

	envSkipped := godotenv.Load(envFile)

	var cfg zap.Config
	if os.Getenv("GENESIS_LOG_MODE") == "production" {
		cfg = zap.NewProductionConfig()
	} else {
		cfg = zap.NewDevelopmentConfig()
	}

	logger, err := cfg.Build(zap.AddCallerSkip(1))
	if err != nil {
		log.Fatal(err)
	}

	if envSkipped != nil {
		logger.Debug(".env file skipped")
	}

	return logger
}()
