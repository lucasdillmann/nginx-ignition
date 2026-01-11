package i18n

import (
	"golang.org/x/text/language"

	"dillmann.com.br/nginx-ignition/core/i18n"
)

func newDictionary() i18n.Dictionary {
	return i18n.Dictionary{
		Language:  language.AmericanEnglish,
		Templates: map[string]string{"key": "value"},
	}
}

func newDictionaryDTO() dictionaryDTO {
	return dictionaryDTO{
		Language:  "en-US",
		Templates: map[string]string{"key": "value"},
	}
}
