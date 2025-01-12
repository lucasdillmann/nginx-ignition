package log

import (
	"fmt"
	"go.uber.org/zap"
)

var (
	logger *zap.Logger
)

func Init() error {
	config := zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:      false,
		Encoding:         "console",
		EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	config.EncoderConfig.EncodeCaller = nil
	loggerInstance, err := config.Build()
	logger = loggerInstance

	return err
}

func Info(message string, values ...any) {
	logger.Info(fmt.Sprintf(message, values...))
}

func Warn(message string, values ...any) {
	logger.Warn(fmt.Sprintf(message, values...))
}

func Error(message string, values ...any) {
	logger.Error(fmt.Sprintf(message, values...))
}

func Fatal(message string, values ...any) {
	logger.Fatal(fmt.Sprintf(message, values...))
}
