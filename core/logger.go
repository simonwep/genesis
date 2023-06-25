package core

import (
	"go.uber.org/zap"
	"log"
)

var Logger *zap.Logger

func init() {
	logger, err := zap.NewProduction()

	if err != nil {
		log.Fatal(err)
	} else {
		Logger = logger
	}
}
