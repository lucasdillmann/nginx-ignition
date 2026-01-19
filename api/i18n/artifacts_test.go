package i18n

import (
	"dillmann.com.br/nginx-ignition/i18n"
)

func newDictionary() i18n.Dictionary {
	return i18n.En_US()
}

func newDictionaryDTO() dictionaryDTO {
	baseDict := i18n.En_US()
	return dictionaryDTO{
		Language: baseDict.Language().String(),
		Messages: baseDict.Raw(),
	}
}
