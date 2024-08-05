package icode

import "go.uber.org/zap"

var logger *zap.Logger

func InitLogger() {
	logger, _ = zap.NewDevelopment()
}
