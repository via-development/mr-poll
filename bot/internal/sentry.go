package internal

import (
	"github.com/getsentry/sentry-go"
	"github.com/via-development/mr-poll/bot/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// TODO:

func NewSentry(config *config.Config, log *zap.Logger) error {
	if config.SentryDSN == "" {
		return nil
	}

	//sentryOptions := zap.
	err := sentry.Init(sentry.ClientOptions{Dsn: config.SentryDSN})
	if err != nil {
		return err
	}

	return nil
}

func zapLevelToSentryLevel(level zapcore.Level) sentry.Level {
	switch level {
	case zapcore.DebugLevel:
		return sentry.LevelDebug
	case zapcore.InfoLevel:
		return sentry.LevelInfo
	case zapcore.WarnLevel:
		return sentry.LevelWarning
	case zapcore.ErrorLevel:
		return sentry.LevelError
	case zapcore.DPanicLevel, zapcore.PanicLevel, zapcore.FatalLevel:
		return sentry.LevelFatal
	default:
		return sentry.LevelError
	}
}
