package logger

import "go.uber.org/zap"

func NewLogger(env string) *zap.Logger {
	var logger *zap.Logger
	var err error

	switch env {
	case "local":
		config := zap.NewDevelopmentConfig()
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		logger, err = config.Build()

	case "prod":
		config := zap.NewProductionConfig()
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
		logger, err = config.Build()

	default:
		// По умолчанию — production
		config := zap.NewProductionConfig()
		logger, err = config.Build()
	}

	if err != nil {
		return nil
	}

	logger = logger.With(zap.String("service", "MainService"))

	return logger
}
