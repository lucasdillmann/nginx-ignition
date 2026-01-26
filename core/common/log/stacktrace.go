package log

import (
	"go.uber.org/zap/zapcore"
)

type stacktrace struct {
	enabled bool
}

func (s *stacktrace) Enabled(_ zapcore.Level) bool {
	return s.enabled
}
