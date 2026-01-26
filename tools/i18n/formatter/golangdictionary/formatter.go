package golangdictionary

import (
	"dillmann.com.br/nginx-ignition/tools/i18n/formatter"
)

const template = `package i18n

import (
	"golang.org/x/text/language"
)

func {{.NormalizedLanguageTag}}() Dictionary {
	return newDictionary(
		language.Make("{{.LanguageTag}}"),
		map[string]string{
{{- range .Messages }}
			Keys.{{.CamelCaseKey}}: "{{.Value}}",
{{- end }}
		},
	)
}
`

func New() formatter.Formatter {
	return formatter.FromTemplate(template)
}
