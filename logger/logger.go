package logger

import (
	"go.uber.org/zap"

	"github.com/echoturing/luxury/conf"
)

var (
	logger      *zap.Logger
	SugarLogger *zap.SugaredLogger
)

func Init(env conf.Env) {
	switch env {
	default:
		fallthrough
	case conf.EnvTest:
		logger, _ = zap.NewDevelopment()
		SugarLogger = logger.Sugar()
	case conf.EnvProd:
		logger, _ = zap.NewProduction()
		SugarLogger = logger.Sugar()
	}
}

func Sync() {
	err := logger.Sync()
	if err != nil {
		SugarLogger.Errorw("sync log error", "err", err.Error())
	}
}
