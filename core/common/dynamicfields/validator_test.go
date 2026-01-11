package dynamicfields

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"dillmann.com.br/nginx-ignition/core/common/validation"
)

func Test_Validate(t *testing.T) {
	t.Run("reports missing required field", func(t *testing.T) {
		dynamicField := newDynamicField()
		dynamicField.ID = "field1"
		dynamicField.Type = SingleLineTextType
		dynamicField.Required = true

		dynamicFields := []DynamicField{*dynamicField}
		err := Validate(t.Context(), dynamicFields, map[string]any{})

		var consistencyErr *validation.ConsistencyError
		assert.ErrorAs(t, err, &consistencyErr)
		assert.Equal(t, "parameters.field1", consistencyErr.Violations[0].Path)
	})

	t.Run("reports invalid email", func(t *testing.T) {
		dynamicField := newDynamicField()
		dynamicField.ID = "email"
		dynamicField.Type = EmailType
		dynamicField.Required = true

		dynamicFields := []DynamicField{*dynamicField}
		assert.Error(t, Validate(t.Context(), dynamicFields, map[string]any{"email": "invalid"}))
	})

	t.Run("reports invalid boolean", func(t *testing.T) {
		dynamicField := newDynamicField()
		dynamicField.ID = "bool"
		dynamicField.Type = BooleanType
		dynamicField.Required = true

		dynamicFields := []DynamicField{*dynamicField}
		assert.Error(t, Validate(t.Context(), dynamicFields, map[string]any{"bool": "not bool"}))
	})

	t.Run("reports invalid enum", func(t *testing.T) {
		dynamicField := newDynamicField()
		dynamicField.ID = "enum"
		dynamicField.Type = EnumType
		dynamicField.Required = true
		dynamicField.EnumOptions = []EnumOption{
			{ID: "opt1"},
		}

		dynamicFields := []DynamicField{*dynamicField}
		assert.Error(t, Validate(t.Context(), dynamicFields, map[string]any{"enum": "invalid"}))
	})
}
