package i18n

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"dillmann.com.br/nginx-ignition/core/i18n"
)

func Test_converter(t *testing.T) {
	t.Run("toDTO", func(t *testing.T) {
		t.Run("maps dictionaries correctly", func(t *testing.T) {
			input := []i18n.Dictionary{newDictionary()}
			expected := []dictionaryDTO{newDictionaryDTO()}

			result := toDTO(input)

			assert.Equal(t, expected, result)
		})

		t.Run("returns empty slice for empty input", func(t *testing.T) {
			result := toDTO([]i18n.Dictionary{})
			assert.Empty(t, result)
		})
	})
}
