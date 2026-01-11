package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

func Test_ConsistencyValidator(t *testing.T) {
	t.Run("Result", func(t *testing.T) {
		t.Run("returns nil when no violations", func(t *testing.T) {
			validator := NewValidator()

			assert.NoError(t, validator.Result())
		})

		t.Run("returns error with violations when added", func(t *testing.T) {
			validator := NewValidator()
			msg1 := i18n.Raw("error message 1")
			msg2 := i18n.Raw("error message 2")
			validator.Add("field1", msg1)
			validator.Add("field2", msg2)

			err := validator.Result()
			assert.Error(t, err)
			var consistencyErr *ConsistencyError
			assert.ErrorAs(t, err, &consistencyErr)
			assert.Len(t, consistencyErr.Violations, 2)
			assert.Equal(t, "field1", consistencyErr.Violations[0].Path)
			assert.Equal(t, msg1, consistencyErr.Violations[0].Message)
			assert.Equal(t, "field2", consistencyErr.Violations[1].Path)
			assert.Equal(t, msg2, consistencyErr.Violations[1].Message)
		})
	})
}
