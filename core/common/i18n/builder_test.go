package i18n

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_builder(t *testing.T) {
	t.Run("M", func(t *testing.T) {
		t.Run("initializes message correctly", func(t *testing.T) {
			key := "some-key"

			message := M(t.Context(), key)

			assert.NotNil(t, message)
			assert.Equal(t, t.Context(), message.ctx)
			assert.Equal(t, key, message.Key)
			assert.NotNil(t, message.Variables)
			assert.Empty(t, message.Variables)
		})
	})

	t.Run("V", func(t *testing.T) {
		t.Run("adds variable and supports chaining", func(t *testing.T) {
			message := M(t.Context(), "key").
				V("var1", "value1").
				V("var2", 123)

			assert.Equal(t, "value1", message.Variables["var1"])
			assert.Equal(t, 123, message.Variables["var2"])
		})

		t.Run("updates existing variable", func(t *testing.T) {
			message := M(t.Context(), "key").
				V("var1", "value1").
				V("var1", "value2")

			assert.Equal(t, "value2", message.Variables["var1"])
		})
	})
}
