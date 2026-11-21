package progress

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_Report(t *testing.T) {
	t.Run("does nothing when there is no reporter", func(t *testing.T) {
		assert.NotPanics(t, func() {
			Report(context.Background(), LevelInfo, "message")
		})
	})

	t.Run("forwards the event to the context reporter", func(t *testing.T) {
		controller := gomock.NewController(t)
		defer controller.Finish()

		var captured Event
		reporter := NewMockedReporter(controller)
		reporter.EXPECT().
			Report(gomock.Any()).
			Do(func(event Event) {
				captured = event
			})

		ctx := WithReporter(context.Background(), reporter)
		Report(ctx, LevelWarn, "message")

		assert.Equal(t, LevelWarn, captured.Level)
		assert.Equal(t, "message", captured.Message)
		assert.False(t, captured.Time.IsZero())
	})
}
