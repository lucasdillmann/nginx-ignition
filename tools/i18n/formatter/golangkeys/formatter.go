package golangkeys

import (
	"dillmann.com.br/nginx-ignition/tools/i18n/formatter"
)

const template = `package i18n

var Keys = struct {
{{- range .Messages }}
	{{.CamelCaseKey}} string
{{- end }}
} {
{{- range .Messages }}
	{{.CamelCaseKey}}: "{{.PropertiesKey}}",
{{- end }}
}
`

func New() formatter.Formatter {
	return formatter.FromTemplate(template)
}
