package logger

import (
	"github.com/TheZeroSlave/zapsentry"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func ConfigureLogger(level, sentryDSN string) (*zap.Logger, error) {
	if level == "" {
		level = "INFO"
	}
	lvl := zapcore.Level(0)
	if err := lvl.UnmarshalText([]byte(level)); err != nil {
		return nil, err
	}
	logger, err := zap.Config{
		Level: zap.NewAtomicLevelAt(lvl),
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
		Encoding:         "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "timestamp",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "message",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
	}.Build()
	if err != nil {
		return nil, err
	}

	if sentryDSN != "" {
		logger, err = addSentryLogger(logger, sentryDSN)
		if err != nil {
			return nil, err
		}
	}

	return logger, nil
}

func addSentryLogger(logger *zap.Logger, dsn string) (*zap.Logger, error) {
	core, err := zapsentry.NewCore(zapsentry.Configuration{
		Level:             zapcore.ErrorLevel,
		EnableBreadcrumbs: true,
		BreadcrumbLevel:   zapcore.ErrorLevel,
		Tags: map[string]string{
			"component": "system",
		},
	}, zapsentry.NewSentryClientFromDSN(dsn))
	if err != nil {
		return nil, err
	}

	logger = logger.With(zapsentry.NewScope())
	return zapsentry.AttachCoreToLogger(core, logger), nil
}
