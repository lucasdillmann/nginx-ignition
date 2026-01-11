package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

func Test_ConsistencyError(t *testing.T) {
	t.Run("Error", func(t *testing.T) {
		t.Run("returns consistent error message", func(t *testing.T) {
			err := NewError([]ConsistencyViolation{
				{
					Path:    "field1",
					Message: i18n.Raw("error 1"),
				},
			})

			assert.Equal(t, "One or more problems where found in the provided value", err.Error())
		})
	})
}

func Test_SingleFieldError(t *testing.T) {
	t.Run("creates error with single violation", func(t *testing.T) {
		message := i18n.Raw("error message")
		err := SingleFieldError("field1", message)

		assert.Len(t, err.Violations, 1)
		assert.Equal(t, "field1", err.Violations[0].Path)
		assert.Equal(t, message, err.Violations[0].Message)
	})
}
