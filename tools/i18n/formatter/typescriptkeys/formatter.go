package typescriptkeys

import (
	"dillmann.com.br/nginx-ignition/tools/i18n/formatter"
)

const template = `enum MessageKey {
{{- range .Messages }}
	{{.CamelCaseKey}} = "{{.PropertiesKey}}",
{{- end }}
}

export default MessageKey
`

func New() formatter.Formatter {
	return formatter.FromTemplate(template)
}
