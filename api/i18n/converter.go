package i18n

import (
	"golang.org/x/text/language"

	"dillmann.com.br/nginx-ignition/i18n"
)

func toDictionaryDTO(dictionary i18n.Dictionary) dictionaryResponseDTO {
	return dictionaryResponseDTO{
		Messages: dictionary.Raw(),
		Language: dictionary.Language().String(),
	}
}

func toAvailableLanguagesDTO(
	dictionaries []i18n.Dictionary,
	defaultLanguage language.Tag,
) availableLanguagesResponseDTO {
	languages := make([]string, len(dictionaries))

	for index := range dictionaries {
		languages[index] = dictionaries[index].Language().String()
	}

	return availableLanguagesResponseDTO{
		Available:       languages,
		DefaultLanguage: defaultLanguage.String(),
	}
}
