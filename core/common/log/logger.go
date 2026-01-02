package log

import (
	"fmt"
	"log"

	"go.uber.org/zap"
)

var delegate *zap.Logger

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
	delegate = loggerInstance

	return err
}

func Info(message string) {
	delegate.Info(message)
}

func Infof(message string, values ...any) {
	delegate.Info(fmt.Sprintf(message, values...))
}

func Warn(message string) {
	delegate.Warn(message)
}

func Warnf(message string, values ...any) {
	delegate.Warn(fmt.Sprintf(message, values...))
}

func Error(message string) {
	delegate.Error(message)
}

func Errorf(message string, values ...any) {
	delegate.Error(fmt.Sprintf(message, values...))
}

func Fatal(message string) {
	delegate.Fatal(message)
}

func Fatalf(message string, values ...any) {
	delegate.Fatal(fmt.Sprintf(message, values...))
}

func Std() *log.Logger {
	return zap.NewStdLog(delegate)
}
