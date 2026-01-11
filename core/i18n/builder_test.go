package i18n

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_builder(t *testing.T) {
	t.Run("K", func(t *testing.T) {
		t.Run("initializes message correctly", func(t *testing.T) {
			ctx := context.Background()
			key := "some-key"

			message := K(ctx, key)

			assert.NotNil(t, message)
			assert.Equal(t, ctx, message.ctx)
			assert.Equal(t, key, message.Key)
			assert.NotNil(t, message.Variables)
			assert.Empty(t, message.Variables)
		})
	})

	t.Run("V", func(t *testing.T) {
		t.Run("adds variable and supports chaining", func(t *testing.T) {
			ctx := context.Background()
			message := K(ctx, "key").
				V("var1", "value1").
				V("var2", 123)

			assert.Equal(t, "value1", message.Variables["var1"])
			assert.Equal(t, 123, message.Variables["var2"])
		})

		t.Run("updates existing variable", func(t *testing.T) {
			ctx := context.Background()
			message := K(ctx, "key").
				V("var1", "value1").
				V("var1", "value2")

			assert.Equal(t, "value2", message.Variables["var1"])
		})
	})
}
