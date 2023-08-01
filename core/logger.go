package core

import (
	"go.uber.org/zap"
	"log"
	"os"
)

var Logger = func() *zap.Logger {
	var logger *zap.Logger
	var err error

	loadEnvFile()
	if os.Getenv("GENESIS_LOG_MODE") == "development" {
		logger, err = zap.NewDevelopmentConfig().Build(
			zap.AddCallerSkip(1),
			zap.IncreaseLevel(zap.DebugLevel),
		)
	} else {
		logger, err = zap.NewProductionConfig().Build(zap.AddCallerSkip(1))
	}

	if err != nil {
		log.Fatal(err)
	}

	return logger
}()
