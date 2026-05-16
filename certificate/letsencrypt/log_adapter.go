package letsencrypt

import (
	"context"
	"log/slog"
	"strings"

	applog "dillmann.com.br/nginx-ignition/core/common/log"
)

type logAdapter struct {
	prependAttrs []slog.Attr
}

func (l *logAdapter) Enabled(context.Context, slog.Level) bool {
	return true
}

func (l *logAdapter) Handle(_ context.Context, record slog.Record) error {
	builder := strings.Builder{}
	_, _ = builder.WriteString(record.Message)

	writeAttr := func(attr slog.Attr) {
		_ = builder.WriteByte(' ')
		_, _ = builder.WriteString(attr.Key)
		_ = builder.WriteByte('=')

		_, _ = builder.WriteString(attr.Value.String())
	}

	for _, attr := range l.prependAttrs {
		writeAttr(attr)
	}

	record.Attrs(func(attribute slog.Attr) bool {
		writeAttr(attribute)
		return true
	})

	message := builder.String()

	switch record.Level {
	case slog.LevelWarn:
		applog.Warn(message)
	case slog.LevelError:
		applog.Error(message)
	case slog.LevelDebug:
		applog.Debug(message)
	default:
		applog.Info(message)
	}

	return nil
}

func (l *logAdapter) WithAttrs(attributes []slog.Attr) slog.Handler {
	return &logAdapter{
		prependAttrs: append(
			append([]slog.Attr(nil), l.prependAttrs...),
			attributes...,
		),
	}
}

func (l *logAdapter) WithGroup(string) slog.Handler {
	return l
}
