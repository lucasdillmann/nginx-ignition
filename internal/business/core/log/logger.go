package log

import (
	"fmt"
	"log"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var delegate *zap.Logger

func init() {
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeCaller = nil

	delegate = zap.New(
		platformCore(zapcore.NewConsoleEncoder(encoderConfig)),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel),
		zap.ErrorOutput(zapcore.Lock(os.Stderr)),
	)
}

func Debug(message string) {
	delegate.Debug(message)
}

func Debugf(message string, values ...any) {
	delegate.Debug(fmt.Sprintf(message, values...))
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
