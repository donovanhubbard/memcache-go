package utils

import (
	"go.uber.org/zap"
)

var Logger *zap.Logger
var Sugar *zap.SugaredLogger

func InitializeLogger() {
	Logger, _ = zap.NewProduction()
	defer Logger.Sync() // flushes buffer, if any
	Sugar = Logger.Sugar()

	Sugar.Info("Logging initialized")
}