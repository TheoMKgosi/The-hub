package config

import (
	"context"

	"github.com/getsentry/sentry-go"
	sentryzap "github.com/getsentry/sentry-go/zap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func InitLogger(dsn, env string, debug bool) error {
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              dsn,
		Environment:      env,
		Debug:            debug,
		AttachStacktrace: true,
		TracesSampleRate: 0.2,
	}); err != nil {
		return err
	}

	ctx := context.Background()
	sentryCore := sentryzap.NewSentryCore(ctx, sentryzap.Option{
		Level: []zapcore.Level{
			zapcore.InfoLevel,
			zapcore.WarnLevel,
			zapcore.ErrorLevel,
		},
		AddCaller: true,
	})
	
	Logger = zap.New(sentryCore)
	return nil
}
