package dynamicfields

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveSensitiveFields(t *testing.T) {
	t.Run("removes sensitive fields from map", func(t *testing.T) {
		values := map[string]any{
			"field1": "value1",
			"field2": "value2",
			"field3": "value3",
		}
		fields := []DynamicField{
			{
				ID:        "field1",
				Sensitive: false,
			},
			{
				ID:        "field2",
				Sensitive: true,
			},
			{
				ID:        "field3",
				Sensitive: false,
			},
		}

		RemoveSensitiveFields(&values, fields)

		assert.Contains(t, values, "field1")
		assert.NotContains(t, values, "field2")
		assert.Contains(t, values, "field3")
	})

	t.Run("handles empty map", func(t *testing.T) {
		values := map[string]any{}
		fields := []DynamicField{
			{
				ID:        "field1",
				Sensitive: true,
			},
		}

		RemoveSensitiveFields(&values, fields)

		assert.Empty(t, values)
	})

	t.Run("handles empty fields list", func(t *testing.T) {
		values := map[string]any{
			"field1": "value1",
		}
		fields := []DynamicField{}

		RemoveSensitiveFields(&values, fields)

		assert.Contains(t, values, "field1")
	})

	t.Run("removes multiple sensitive fields", func(t *testing.T) {
		values := map[string]any{
			"field1": "value1",
			"field2": "value2",
			"field3": "value3",
			"field4": "value4",
		}
		fields := []DynamicField{
			{
				ID:        "field1",
				Sensitive: true,
			},
			{
				ID:        "field2",
				Sensitive: false,
			},
			{
				ID:        "field3",
				Sensitive: true,
			},
			{
				ID:        "field4",
				Sensitive: true,
			},
		}

		RemoveSensitiveFields(&values, fields)

		assert.NotContains(t, values, "field1")
		assert.Contains(t, values, "field2")
		assert.NotContains(t, values, "field3")
		assert.NotContains(t, values, "field4")
	})

	t.Run("do nothing for non-existent fields", func(t *testing.T) {
		values := map[string]any{
			"field1": "value1",
		}
		fields := []DynamicField{
			{
				ID:        "nonexistent",
				Sensitive: true,
			},
		}

		RemoveSensitiveFields(&values, fields)

		assert.Contains(t, values, "field1")
	})
}
