package i18n

import (
	"golang.org/x/text/language"

	"dillmann.com.br/nginx-ignition/i18n/dict"
)

func toDTO(
	defaultLanguage language.Tag,
	dictionaries []dict.Dictionary,
) dictionariesDTO {
	output := make([]dictionaryDTO, len(dictionaries))

	for index := range dictionaries {
		output[index] = dictionaryDTO{
			Language:  dictionaries[index].Language().String(),
			Templates: dictionaries[index].Templates(),
		}
	}

	return dictionariesDTO{
		DefaultLanguage: defaultLanguage.String(),
		Dictionaries:    output,
	}
}
