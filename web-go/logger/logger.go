package logger

import "go.uber.org/zap"

func Setup() {
	logger, _ := zap.NewProduction()
	zap.ReplaceGlobals(logger)
}
