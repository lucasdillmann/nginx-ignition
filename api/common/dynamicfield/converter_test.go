package dynamicfield

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

func Test_toResponse(t *testing.T) {
	t.Run("converts domain model to DTO", func(t *testing.T) {
		input := []dynamicfields.DynamicField{
			{
				ID:          "field1",
				Type:        dynamicfields.SingleLineTextType,
				Required:    true,
				Priority:    1,
				Description: "Desc 1",
			},
			{
				ID:           "field2",
				Type:         dynamicfields.BooleanType,
				DefaultValue: "10",
			},
			{
				ID:   "field3",
				Type: dynamicfields.EnumType,
				EnumOptions: []dynamicfields.EnumOption{
					{
						ID:          "A",
						Description: "Desc A",
					},
				},
			},
		}

		result := ToResponse(input)

		assert.Len(t, result, 3)
		assert.Equal(t, "field1", result[0].ID)
		assert.Equal(t, "SINGLE_LINE_TEXT", result[0].Type)
		assert.True(t, result[0].Required)
		assert.Equal(t, 1, result[0].Priority)
		assert.Equal(t, "Desc 1", result[0].Description)

		assert.Equal(t, "field2", result[1].ID)
		assert.Equal(t, "BOOLEAN", result[1].Type)
		assert.Equal(t, "10", result[1].DefaultValue)

		assert.Equal(t, "field3", result[2].ID)
		assert.Equal(t, "ENUM", result[2].Type)
		assert.Len(t, result[2].EnumOptions, 1)
		assert.Equal(t, "A", result[2].EnumOptions[0].ID)
	})

	t.Run("returns empty slice when input is nil or empty", func(t *testing.T) {
		assert.Empty(t, ToResponse(nil))
		assert.Empty(t, ToResponse([]dynamicfields.DynamicField{}))
	})
}
