package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ConsistencyValidator_Result(t *testing.T) {
	t.Run("returns nil when no violations", func(t *testing.T) {
		val := NewValidator()

		assert.NoError(t, val.Result())
	})

	t.Run("returns error with violations when added", func(t *testing.T) {
		val := NewValidator()
		val.Add("field1", "error message 1")
		val.Add("field2", "error message 2")

		err := val.Result()
		assert.Error(t, err)
		var consistencyErr *ConsistencyError
		assert.ErrorAs(t, err, &consistencyErr)
		assert.Len(t, consistencyErr.Violations, 2)
		assert.Equal(t, "field1", consistencyErr.Violations[0].Path)
		assert.Equal(t, "error message 1", consistencyErr.Violations[0].Message)
		assert.Equal(t, "field2", consistencyErr.Violations[1].Path)
		assert.Equal(t, "error message 2", consistencyErr.Violations[1].Message)
	})
}
