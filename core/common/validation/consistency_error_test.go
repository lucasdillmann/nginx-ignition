package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ConsistencyError_Error(t *testing.T) {
	t.Run("returns consistent error message", func(t *testing.T) {
		err := NewError([]ConsistencyViolation{
			{
				Path:    "field1",
				Message: "error 1",
			},
		})

		assert.Equal(t, "One or more problems where found in the provided value", err.Error())
	})
}

func Test_SingleFieldError(t *testing.T) {
	t.Run("creates error with single violation", func(t *testing.T) {
		err := SingleFieldError("field1", "error message")

		assert.Len(t, err.Violations, 1)
		assert.Equal(t, "field1", err.Violations[0].Path)
		assert.Equal(t, "error message", err.Violations[0].Message)
	})
}
