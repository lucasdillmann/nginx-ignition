package progress

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_ReporterFrom(t *testing.T) {
	t.Run("returns the reporter from the context", func(t *testing.T) {
		controller := gomock.NewController(t)
		defer controller.Finish()

		reporter := NewMockedReporter(controller)
		ctx := WithReporter(context.Background(), reporter)

		assert.Equal(t, reporter, ReporterFrom(ctx))
	})

	t.Run("returns nil when the context has no reporter", func(t *testing.T) {
		assert.Nil(t, ReporterFrom(context.Background()))
	})
}

func Test_WithReporter(t *testing.T) {
	t.Run("ignores a nil reporter", func(t *testing.T) {
		ctx := WithReporter(context.Background(), nil)

		assert.Nil(t, ReporterFrom(ctx))
	})
}
