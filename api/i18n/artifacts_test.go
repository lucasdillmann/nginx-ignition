package i18n

import (
	"dillmann.com.br/nginx-ignition/core/common/i18n/dict"
)

func newDictionary() dict.Dictionary {
	return dict.EnUS()
}

func newDictionaryDTO() dictionaryDTO {
	baseDict := dict.EnUS()
	return dictionaryDTO{
		Language:  baseDict.Language().String(),
		Templates: baseDict.Templates(),
	}
}
