package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// TODO: move to config
const logPath = "./logs/go.log"

func NewLogger(env string) *zap.Logger {
	var cfg zap.Config

	os.OpenFile(logPath, os.O_RDONLY|os.O_CREATE, 0666)

	switch env {
	case "local":
		cfg = zap.NewDevelopmentConfig()
		cfg.Encoding = "json"
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)

	case "prod":
		cfg = zap.NewProductionConfig()
		cfg.Encoding = "json"
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)

	default:
		cfg = zap.NewProductionConfig()
		cfg.Encoding = "json"
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	cfg.OutputPaths = []string{"stdout", logPath}

	logger, err := cfg.Build()
	if err != nil {
		return nil
	}

	logger = logger.With(zap.String("service", "MainService"))

	return logger
}
