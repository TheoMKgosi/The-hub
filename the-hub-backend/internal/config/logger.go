package config

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var Logger *zap.SugaredLogger

func InitLogger() {
	env := os.Getenv("APP_ENV") // e.g. "production" or "development"

	var logger *zap.Logger
	var err error

	if env == "production" {
		logger, err = zap.NewProduction()
	} else {
		// Development config: console encoder + debug level + caller info
		cfg := zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		logger, err = cfg.Build()
	}

	if err != nil {
		panic(err)
	}

	Logger = logger.Sugar()
}

