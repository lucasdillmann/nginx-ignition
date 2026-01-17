package dynamicfields

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_RemoveSensitiveFields(t *testing.T) {
	t.Run("removes sensitive fields from map", func(t *testing.T) {
		values := map[string]any{
			"field1": "value1",
			"field2": "value2",
			"field3": "value3",
		}

		dynamicField1 := newDynamicField(t.Context())
		dynamicField1.ID = "field1"

		dynamicField2 := newDynamicField(t.Context())
		dynamicField2.ID = "field2"
		dynamicField2.Sensitive = true

		dynamicField3 := newDynamicField(t.Context())
		dynamicField3.ID = "field3"

		dynamicFields := []DynamicField{*dynamicField1, *dynamicField2, *dynamicField3}

		RemoveSensitiveFields(&values, dynamicFields)

		assert.Contains(t, values, "field1")
		assert.NotContains(t, values, "field2")
		assert.Contains(t, values, "field3")
	})

	t.Run("handles empty map", func(t *testing.T) {
		values := map[string]any{}

		dynamicField := newDynamicField(t.Context())
		dynamicField.ID = "field1"
		dynamicField.Sensitive = true

		dynamicFields := []DynamicField{*dynamicField}

		RemoveSensitiveFields(&values, dynamicFields)

		assert.Empty(t, values)
	})

	t.Run("handles empty fields list", func(t *testing.T) {
		values := map[string]any{
			"field1": "value1",
		}
		dynamicFields := []DynamicField{}

		RemoveSensitiveFields(&values, dynamicFields)

		assert.Contains(t, values, "field1")
	})

	t.Run("removes multiple sensitive fields", func(t *testing.T) {
		values := map[string]any{
			"field1": "value1",
			"field2": "value2",
			"field3": "value3",
			"field4": "value4",
		}

		dynamicField1 := newDynamicField(t.Context())
		dynamicField1.ID = "field1"
		dynamicField1.Sensitive = true

		dynamicField2 := newDynamicField(t.Context())
		dynamicField2.ID = "field2"

		dynamicField3 := newDynamicField(t.Context())
		dynamicField3.ID = "field3"
		dynamicField3.Sensitive = true

		dynamicField4 := newDynamicField(t.Context())
		dynamicField4.ID = "field4"
		dynamicField4.Sensitive = true

		dynamicFields := []DynamicField{
			*dynamicField1,
			*dynamicField2,
			*dynamicField3,
			*dynamicField4,
		}

		RemoveSensitiveFields(&values, dynamicFields)

		assert.NotContains(t, values, "field1")
		assert.Contains(t, values, "field2")
		assert.NotContains(t, values, "field3")
		assert.NotContains(t, values, "field4")
	})

	t.Run("do nothing for non-existent fields", func(t *testing.T) {
		values := map[string]any{
			"field1": "value1",
		}

		dynamicField := newDynamicField(t.Context())
		dynamicField.ID = "nonexistent"
		dynamicField.Sensitive = true

		dynamicFields := []DynamicField{*dynamicField}

		RemoveSensitiveFields(&values, dynamicFields)

		assert.Contains(t, values, "field1")
	})
}
