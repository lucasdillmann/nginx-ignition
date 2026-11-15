//go:build !windows

package log

import (
	"os"

	"go.uber.org/zap/zapcore"
)

func platformCore(encoder zapcore.Encoder) zapcore.Core {
	return zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), zapcore.InfoLevel)
}
