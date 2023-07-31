package core

import (
	"go.uber.org/zap"
	"log"
)

var logger = func() *zap.Logger {
	var logger *zap.Logger
	var err error

	if Config.AppLogMode == "development" {
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

func GetLogger() *zap.Logger {
	return logger
}

func Debug(message string, fields ...zap.Field) {
	logger.Debug(message, fields...)
}

func Warn(message string, fields ...zap.Field) {
	logger.Warn(message, fields...)
}

func Error(message string, fields ...zap.Field) {
	logger.Error(message, fields...)
}

func Fatal(message string, fields ...zap.Field) {
	logger.Fatal(message, fields...)
}
