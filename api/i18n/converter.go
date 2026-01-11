package i18n

import (
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

func toDTO(dictionaries []i18n.Dictionary) []dictionaryDTO {
	output := make([]dictionaryDTO, len(dictionaries))

	for index, dictionary := range dictionaries {
		output[index] = dictionaryDTO{
			Language:  dictionary.Language.String(),
			Templates: dictionary.Templates,
		}
	}

	return output
}
