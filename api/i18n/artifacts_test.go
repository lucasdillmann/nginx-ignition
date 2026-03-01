package i18n

import (
	"dillmann.com.br/nginx-ignition/i18n"
)

func newDictionary() i18n.Dictionary {
	return i18n.En()
}

func newDictionaryDTO() dictionaryResponseDTO {
	baseDict := i18n.En()
	return dictionaryResponseDTO{
		Language: baseDict.Language().String(),
		Messages: baseDict.Raw(),
	}
}
