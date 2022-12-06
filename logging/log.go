package logging

import "go.uber.org/zap"

var Logger *zap.Logger

var SugaredLogger *zap.SugaredLogger

func init() {

	zapLogger, _ := zap.NewDevelopment()

	Logger = zapLogger

	SugaredLogger = zapLogger.Sugar()

}
