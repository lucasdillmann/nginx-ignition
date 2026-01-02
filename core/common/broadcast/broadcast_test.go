package broadcast

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"dillmann.com.br/nginx-ignition/core/common/log"
)

func init() {
	_ = log.Init()
}

func Test_Listen(t *testing.T) {
	t.Run("creates channel for new qualifier", func(t *testing.T) {
		qualifier := "test-qualifier-1"
		ch := Listen(qualifier)

		assert.NotNil(t, ch)
	})

	t.Run("returns same channel for existing qualifier", func(t *testing.T) {
		qualifier := "test-qualifier-2"
		ch1 := Listen(qualifier)
		ch2 := Listen(qualifier)

		assert.Equal(t, ch1, ch2)
	})
}

func Test_SendSignal(t *testing.T) {
	t.Run("sends context to listening channel", func(t *testing.T) {
		qualifier := "test-qualifier-3"
		ch := Listen(qualifier)

		ctx := context.Background()
		go SendSignal(ctx, qualifier)

		select {
		case received := <-ch:
			assert.Equal(t, ctx, received)
		case <-time.After(time.Second):
			t.Fatal("signal not received")
		}
	})

	t.Run("does not block when no listeners", func(_ *testing.T) {
		ctx := context.Background()
		SendSignal(ctx, "non-existent-qualifier-2")
	})
}

func Test_Shutdown_Run(t *testing.T) {
	t.Run("closes all channels", func(t *testing.T) {
		ch1 := Listen("test-shutdown-1")
		ch2 := Listen("test-shutdown-2")

		s := shutdown{}
		s.Run(context.Background())

		_, ok1 := <-ch1
		_, ok2 := <-ch2

		assert.False(t, ok1)
		assert.False(t, ok2)
	})
}
