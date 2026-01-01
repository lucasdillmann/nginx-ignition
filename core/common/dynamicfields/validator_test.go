package dynamicfields

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"dillmann.com.br/nginx-ignition/core/common/validation"
)

func TestValidate(t *testing.T) {
	t.Run("reports missing required field", func(t *testing.T) {
		fields := []DynamicField{
			{
				ID:       "field1",
				Type:     SingleLineTextType,
				Required: true,
			},
		}
		err := Validate(fields, map[string]any{})

		var consistencyErr *validation.ConsistencyError
		assert.ErrorAs(t, err, &consistencyErr)
		assert.Equal(t, "parameters.field1", consistencyErr.Violations[0].Path)
	})

	t.Run("reports invalid email", func(t *testing.T) {
		fields := []DynamicField{
			{
				ID:       "email",
				Type:     EmailType,
				Required: true,
			},
		}
		assert.Error(t, Validate(fields, map[string]any{"email": "invalid"}))
	})

	t.Run("reports invalid boolean", func(t *testing.T) {
		fields := []DynamicField{
			{
				ID:       "bool",
				Type:     BooleanType,
				Required: true,
			},
		}
		assert.Error(t, Validate(fields, map[string]any{"bool": "not bool"}))
	})

	t.Run("reports invalid enum", func(t *testing.T) {
		fields := []DynamicField{
			{
				ID:       "enum",
				Type:     EnumType,
				Required: true,
				EnumOptions: []EnumOption{
					{ID: "opt1"},
				},
			},
		}
		assert.Error(t, Validate(fields, map[string]any{"enum": "invalid"}))
	})
}
