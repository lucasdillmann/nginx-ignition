package letsencrypt

import (
	"context"
	"log/slog"
	"strings"

	applog "github.com/lucasdillmann/nginx-ignition/internal/business/core/log"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/progress"
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

	switch {
	case record.Level >= slog.LevelError:
		applog.Error(message)
	case record.Level >= slog.LevelWarn:
		applog.Warn(message)
	case record.Level <= slog.LevelDebug:
		applog.Debug(message)
	default:
		applog.Info(message)
	}

	reportActive(toProgressLevel(record.Level), message)

	return nil
}

func toProgressLevel(level slog.Level) progress.Level {
	switch {
	case level >= slog.LevelError:
		return progress.LevelError
	case level >= slog.LevelWarn:
		return progress.LevelWarn
	case level <= slog.LevelDebug:
		return progress.LevelDebug
	default:
		return progress.LevelInfo
	}
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
