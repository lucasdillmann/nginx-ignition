package eventstream

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/progress"
)

func Test_streamReporter(t *testing.T) {
	t.Run("forwards events", func(t *testing.T) {
		events := make(chan progress.Event, 1)
		done := make(chan struct{})
		reporter := &streamReporter{events: events, done: done}
		event := progress.Event{Level: progress.LevelInfo, Message: "message"}

		reporter.Report(event)

		assert.Equal(t, event, <-events)
	})

	t.Run("unblocks when the stream is done", func(t *testing.T) {
		events := make(chan progress.Event, 1)
		events <- progress.Event{}
		done := make(chan struct{})
		close(done)
		reporter := &streamReporter{events: events, done: done}

		completed := make(chan struct{})
		go func() {
			reporter.Report(progress.Event{})
			close(completed)
		}()

		select {
		case <-completed:
		case <-time.After(time.Second):
			t.Fatal("reporter did not unblock after the stream was done")
		}
	})
}
