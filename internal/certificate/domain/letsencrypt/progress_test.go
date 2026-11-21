package letsencrypt

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_withActiveReporter(t *testing.T) {
	t.Run("returns nil when there is no reporter", func(t *testing.T) {
		assert.Nil(t, withActiveReporter(context.Background()))
	})
}
