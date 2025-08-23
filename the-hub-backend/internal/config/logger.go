package config

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

var Logger *zap.SugaredLogger

func InitLogger() {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	var logger *zap.Logger
	var err error

	if env == "production" {
		// Log rotation using lumberjack
		lumberjackLogger := &lumberjack.Logger{
			Filename:   "logs/app.log",
			MaxSize:    50, // MB
			MaxBackups: 7,
			MaxAge:     30,   // days
			Compress:   true, // gzip old files
		}

		writeSyncer := zapcore.NewMultiWriteSyncer(
			zapcore.AddSync(os.Stdout),        // optional console logging
			zapcore.AddSync(lumberjackLogger), // file logging
		)

		encoderCfg := zap.NewProductionEncoderConfig()
		encoderCfg.TimeKey = "time"
		encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderCfg),
			writeSyncer,
			zap.InfoLevel,
		)

		logger = zap.New(core, zap.AddCaller())
	} else {
		// Dev: coloured console logs
		cfg := zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		logger, err = cfg.Build(zap.AddCaller())
		if err != nil {
			panic("Failed to initialize logger: " + err.Error())
		}
	}

	Logger = logger.Sugar()
}
