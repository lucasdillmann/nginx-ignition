package progress

import (
	"context"
	"time"
)

func Report(ctx context.Context, level Level, message string) {
	reporter := ReporterFrom(ctx)
	if reporter == nil {
		return
	}

	reporter.Report(Event{
		Time:    time.Now(),
		Level:   level,
		Message: message,
	})
}
