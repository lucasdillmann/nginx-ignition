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

func Infof(message string, values ...any) {
	logger.Info(fmt.Sprintf(message, values...))
}

func Warnf(message string, values ...any) {
	logger.Warn(fmt.Sprintf(message, values...))
}

func Errorf(message string, values ...any) {
	logger.Error(fmt.Sprintf(message, values...))
}

func Fatalf(message string, values ...any) {
	logger.Fatal(fmt.Sprintf(message, values...))
}
