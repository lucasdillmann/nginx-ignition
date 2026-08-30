package i18n

import (
	"github.com/lucasdillmann/nginx-ignition/internal/i18n"
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
