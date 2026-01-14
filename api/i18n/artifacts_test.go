package i18n

import (
	"dillmann.com.br/nginx-ignition/i18n"
)

func newDictionary() i18n.Dictionary {
	return i18n.ENUS()
}

func newDictionaryDTO() dictionaryDTO {
	baseDict := i18n.ENUS()
	return dictionaryDTO{
		Language: baseDict.Language().String(),
		Messages: baseDict.Raw(),
	}
}
