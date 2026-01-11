package i18n

import (
	"golang.org/x/text/language"
)

var ptBR = Dictionary{
	Language: language.Make("pt-BR"),
	Templates: map[string]string{
		"value-missing": "Preenchimento obrigatório",
		"test-msg":      "Apenas um teste: ${test-var}",
	},
}
