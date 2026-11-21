package letsencrypt

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/progress"
)

func Test_logAdapter_Handle(t *testing.T) {
	t.Run("forwards records to the active reporter", func(t *testing.T) {
		controller := gomock.NewController(t)
		defer controller.Finish()

		var captured progress.Event
		reporter := progress.NewMockedReporter(controller)
		reporter.EXPECT().
			Report(gomock.Any()).
			Do(func(event progress.Event) {
				captured = event
			})

		restore := withActiveReporter(progress.WithReporter(context.Background(), reporter))
		defer restore()

		record := slog.NewRecord(time.Now(), slog.LevelWarn, "message", 0)
		record.AddAttrs(slog.String("key", "value"))

		adapter := &logAdapter{}
		require.NoError(t, adapter.Handle(context.Background(), record))

		assert.Equal(t, progress.LevelWarn, captured.Level)
		assert.Contains(t, captured.Message, "message")
		assert.Contains(t, captured.Message, "key=value")
	})

	t.Run("does nothing when there is no active reporter", func(t *testing.T) {
		record := slog.NewRecord(time.Now(), slog.LevelInfo, "message", 0)

		adapter := &logAdapter{}

		assert.NotPanics(t, func() {
			_ = adapter.Handle(context.Background(), record)
		})
	})
}

func Test_toProgressLevel(t *testing.T) {
	assert.Equal(t, progress.LevelDebug, toProgressLevel(slog.LevelDebug))
	assert.Equal(t, progress.LevelInfo, toProgressLevel(slog.LevelInfo))
	assert.Equal(t, progress.LevelWarn, toProgressLevel(slog.LevelWarn))
	assert.Equal(t, progress.LevelError, toProgressLevel(slog.LevelError))
}
