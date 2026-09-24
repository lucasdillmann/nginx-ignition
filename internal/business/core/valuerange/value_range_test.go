package valuerange

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_New(t *testing.T) {
	t.Run("creates value range with min and max", func(t *testing.T) {
		minValue := 1
		maxValue := 100
		vr := New(minValue, maxValue)

		assert.NotNil(t, vr)
		assert.Equal(t, minValue, vr.Min)
		assert.Equal(t, maxValue, vr.Max)
	})
}

func Test_ValueRange(t *testing.T) {
	t.Run("Contains", func(t *testing.T) {
		t.Run("returns true for value within range", func(t *testing.T) {
			vr := New(1, 100)

			assert.True(t, vr.Contains(50))
			assert.True(t, vr.Contains(1))
			assert.True(t, vr.Contains(100))
		})

		t.Run("returns false for value below range", func(t *testing.T) {
			vr := New(1, 100)

			assert.False(t, vr.Contains(0))
			assert.False(t, vr.Contains(-1))
		})

		t.Run("returns false for value above range", func(t *testing.T) {
			vr := New(1, 100)

			assert.False(t, vr.Contains(101))
			assert.False(t, vr.Contains(200))
		})
	})
}
