package converter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_wrap(t *testing.T) {
	t.Run("returns converted value when input is not nil", func(t *testing.T) {
		convertFunc := func(_ *int) string {
			return "value"
		}
		input := 1
		result := Wrap(t.Context(), convertFunc, &input)
		assert.Equal(t, "value", result)
	})

	t.Run("panics when input is nil", func(t *testing.T) {
		convertFunc := func(in *int) string {
			val := *in
			return "value" + string(rune(val))
		}
		assert.Panics(t, func() {
			Wrap(t.Context(), convertFunc, nil)
		})
	})
}

func Test_wrap2(t *testing.T) {
	t.Run("returns converted value when inputs are valid", func(t *testing.T) {
		convertFunc := func(_, _ int) string {
			return "value"
		}
		result := Wrap2(t.Context(), convertFunc, 1, 2)
		assert.Equal(t, "value", result)
	})

	t.Run("panics when converter panics", func(t *testing.T) {
		convertFunc := func(_, _ int) string {
			panic("oops")
		}

		assert.Panics(t, func() {
			Wrap2(t.Context(), convertFunc, 1, 2)
		})
	})
}
