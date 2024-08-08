package web_logger

import "go.uber.org/zap"

var logger *zap.Logger
var Log = logger

func InitLogger() {
	logger, _ = zap.NewDevelopment()
}

func Sync() error {
	return logger.Sync()
}
