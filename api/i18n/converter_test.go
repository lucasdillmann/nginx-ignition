package i18n

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"dillmann.com.br/nginx-ignition/core/common/i18n/dict"
)

func Test_converter(t *testing.T) {
	t.Run("toDTO", func(t *testing.T) {
		t.Run("maps dictionaries correctly", func(t *testing.T) {
			input := []dict.Dictionary{newDictionary()}
			expected := []dictionaryDTO{newDictionaryDTO()}

			result := toDTO(input)

			assert.Equal(t, expected, result)
		})

		t.Run("returns empty slice for empty input", func(t *testing.T) {
			result := toDTO([]dict.Dictionary{})
			assert.Empty(t, result)
		})
	})
}
