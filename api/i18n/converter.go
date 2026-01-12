package i18n

import (
	"dillmann.com.br/nginx-ignition/core/common/i18n/dict"
)

func toDTO(dictionaries []dict.Dictionary) []dictionaryDTO {
	output := make([]dictionaryDTO, len(dictionaries))

	for index := range dictionaries {
		output[index] = dictionaryDTO{
			Language:  dictionaries[index].Language().String(),
			Templates: dictionaries[index].Templates(),
		}
	}

	return output
}
