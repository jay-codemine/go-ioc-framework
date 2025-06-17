package zapmod

import (
	"go.uber.org/zap"
)

var Logger *zap.Logger

func Init() *zap.Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	Logger = logger
	return Logger
}
