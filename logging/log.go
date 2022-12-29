package logging

import (
	"go.uber.org/zap"
)

var Logger *zap.Logger

var SugaredLogger *zap.SugaredLogger

func init() {

	config := zap.NewDevelopmentConfig()

	config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)

	zapLogger, _ := config.Build()

	Logger = zapLogger

	SugaredLogger = zapLogger.Sugar()
}
