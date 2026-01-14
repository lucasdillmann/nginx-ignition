package i18n

import (
	"golang.org/x/text/language"

	"dillmann.com.br/nginx-ignition/i18n"
)

func toDTO(
	defaultLanguage language.Tag,
	dictionaries []i18n.Dictionary,
) dictionariesDTO {
	output := make([]dictionaryDTO, len(dictionaries))

	for index := range dictionaries {
		output[index] = dictionaryDTO{
			Language: dictionaries[index].Language().String(),
			Messages: dictionaries[index].Raw(),
		}
	}

	return dictionariesDTO{
		DefaultLanguage: defaultLanguage.String(),
		Dictionaries:    output,
	}
}
