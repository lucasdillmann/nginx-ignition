package ptr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Of(t *testing.T) {
	t.Run("returns pointer to value", func(t *testing.T) {
		value := "test"
		result := Of(value)

		assert.Equal(t, value, *result)
		assert.NotSame(t, &value, result)
	})
}
