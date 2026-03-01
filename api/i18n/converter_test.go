package i18n

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"

	"dillmann.com.br/nginx-ignition/i18n"
)

func Test_converter(t *testing.T) {
	t.Run("toDictionaryDTO", func(t *testing.T) {
		t.Run("maps dictionary correctly", func(t *testing.T) {
			dict := newDictionary()
			expected := newDictionaryDTO()

			result := toDictionaryDTO(dict)

			assert.Equal(t, expected, result)
		})
	})

	t.Run("toAvailableLanguagesDTO", func(t *testing.T) {
		t.Run("maps available languages correctly", func(t *testing.T) {
			dictionaries := []i18n.Dictionary{newDictionary()}
			defaultLanguage := language.AmericanEnglish
			expected := availableLanguagesResponseDTO{
				DefaultLanguage: defaultLanguage.String(),
				Available:       []string{newDictionary().Language().String()},
			}

			result := toAvailableLanguagesDTO(dictionaries, defaultLanguage)

			assert.Equal(t, expected, result)
		})

		t.Run("returns empty slice for empty input", func(t *testing.T) {
			defaultLanguage := language.AmericanEnglish
			expected := availableLanguagesResponseDTO{
				DefaultLanguage: defaultLanguage.String(),
				Available:       []string{},
			}

			result := toAvailableLanguagesDTO(nil, defaultLanguage)

			assert.Equal(t, expected, result)
		})
	})
}
