//go:build windows

package log

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/constants"
	"go.uber.org/zap/zapcore"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/eventlog"
)

const (
	eventIDInformation = 1
	eventIDWarning     = 2
	eventIDError       = 3

	maxEventLogMessageLength = 30000
)

func platformCore(encoder zapcore.Encoder) zapcore.Core {
	stdoutCore := zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), zapcore.InfoLevel)

	isService, err := svc.IsWindowsService()
	if err != nil || !isService {
		return stdoutCore
	}

	eventLog, err := openEventLog()
	if err != nil {
		return stdoutCore
	}

	return &eventLogCore{
		level: zapcore.InfoLevel,
		log:   eventLog,
	}
}

func openEventLog() (*eventlog.Log, error) {
	eventLog, err := eventlog.Open(constants.ServiceName)
	if err == nil {
		return eventLog, nil
	}

	if installErr := eventlog.InstallAsEventCreate(
		constants.ServiceName,
		eventlog.Info|eventlog.Warning|eventlog.Error,
	); installErr != nil {
		return nil, err
	}

	return eventlog.Open(constants.ServiceName)
}

type eventLogCore struct {
	level  zapcore.LevelEnabler
	log    *eventlog.Log
	fields []zapcore.Field
}

func (c *eventLogCore) Enabled(level zapcore.Level) bool {
	return c.level.Enabled(level)
}

func (c *eventLogCore) With(fields []zapcore.Field) zapcore.Core {
	clone := *c

	clone.fields = make([]zapcore.Field, 0, len(c.fields)+len(fields))
	clone.fields = append(clone.fields, c.fields...)
	clone.fields = append(clone.fields, fields...)

	return &clone
}

func (c *eventLogCore) Check(entry zapcore.Entry, checked *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(entry.Level) {
		return checked.AddCore(entry, c)
	}

	return checked
}

func (c *eventLogCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	allFields := make([]zapcore.Field, 0, len(c.fields)+len(fields))
	allFields = append(allFields, c.fields...)
	allFields = append(allFields, fields...)

	message := truncateMessage(entry.Message + formatFields(allFields))

	switch entry.Level {
	case zapcore.DebugLevel, zapcore.InfoLevel:
		return c.log.Info(eventIDInformation, message)
	case zapcore.WarnLevel:
		return c.log.Warning(eventIDWarning, message)
	default:
		return c.log.Error(eventIDError, message)
	}
}

func (c *eventLogCore) Sync() error {
	return nil
}

func formatFields(fields []zapcore.Field) string {
	if len(fields) == 0 {
		return ""
	}

	encoder := zapcore.NewMapObjectEncoder()
	for _, field := range fields {
		field.AddTo(encoder)
	}

	parts := make([]string, 0, len(encoder.Fields))
	for key, value := range encoder.Fields {
		parts = append(parts, fmt.Sprintf("%s=%v", key, value))
	}

	sort.Strings(parts)

	return " " + strings.Join(parts, " ")
}

func truncateMessage(message string) string {
	const suffix = "..."

	if len(message) <= maxEventLogMessageLength {
		return message
	}

	cut := maxEventLogMessageLength - len(suffix)
	for cut > 0 && !utf8.RuneStart(message[cut]) {
		cut--
	}

	return message[:cut] + suffix
}
