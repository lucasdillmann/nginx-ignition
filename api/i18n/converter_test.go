package i18n

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"

	"dillmann.com.br/nginx-ignition/i18n"
)

func Test_converter(t *testing.T) {
	t.Run("toDTO", func(t *testing.T) {
		t.Run("maps dictionaries correctly", func(t *testing.T) {
			dictionaries := []i18n.Dictionary{newDictionary()}
			defaultLanguage := language.AmericanEnglish
			expected := dictionariesDTO{
				DefaultLanguage: defaultLanguage.String(),
				Dictionaries:    []dictionaryDTO{newDictionaryDTO()},
			}

			result := toDTO(defaultLanguage, dictionaries)

			assert.Equal(t, expected, result)
		})

		t.Run("returns empty slice for empty input", func(t *testing.T) {
			defaultLanguage := language.AmericanEnglish
			expected := dictionariesDTO{
				DefaultLanguage: defaultLanguage.String(),
				Dictionaries:    []dictionaryDTO{},
			}

			result := toDTO(defaultLanguage, nil)

			assert.Equal(t, expected, result)
		})
	})
}
