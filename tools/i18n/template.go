package main

const keyFileTemplate = `
package i18n

var Keys = struct {
{{- range .Items }}
	{{.ConstName}} string
{{- end }}
} {
{{- range .Items }}
	{{.ConstName}}: "{{.Value}}",
{{- end }}
}
`

const dictionaryFileTemplate = `
package i18n

import (
	"golang.org/x/text/language"
)

func {{.FuncName}}() Dictionary {
	return newDictionary(
		language.Make("{{.LangTag}}"),
		map[string]string{
{{- range .Keys }}
			Keys.{{.ConstName}}: "{{.Value}}",
{{- end }}
		},
	)
}
`
