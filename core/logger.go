package core

import (
	"go.uber.org/zap"
	"log"
)

var Logger = func() *zap.Logger {
	logger, err := zap.NewProductionConfig().Build(zap.AddCallerSkip(1))

	if err != nil {
		log.Fatal(err)
	}

	return logger
}()
