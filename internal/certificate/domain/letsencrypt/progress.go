package letsencrypt

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	applog "github.com/lucasdillmann/nginx-ignition/internal/business/core/log"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/progress"
)

type activeReporter struct {
	reporter progress.Reporter
}

var currentReporter atomic.Pointer[activeReporter]

func withActiveReporter(ctx context.Context) func() {
	reporter := progress.ReporterFrom(ctx)
	if reporter == nil {
		return nil
	}

	previous := currentReporter.Swap(&activeReporter{reporter: reporter})

	return func() {
		currentReporter.Store(previous)
	}
}

func reportActive(level progress.Level, message string) {
	reporter := currentReporter.Load()
	if reporter == nil {
		return
	}

	reporter.reporter.Report(progress.Event{
		Time:    time.Now(),
		Level:   level,
		Message: message,
	})
}

func reportInfo(ctx context.Context, format string, values ...any) {
	report(ctx, progress.LevelInfo, format, values...)
}

func report(ctx context.Context, level progress.Level, format string, values ...any) {
	message := format
	if len(values) > 0 {
		message = fmt.Sprintf(format, values...)
	}

	switch level {
	case progress.LevelDebug:
		applog.Debug(message)
	case progress.LevelWarn:
		applog.Warn(message)
	case progress.LevelError:
		applog.Error(message)
	default:
		applog.Info(message)
	}

	progress.Report(ctx, level, message)
}
