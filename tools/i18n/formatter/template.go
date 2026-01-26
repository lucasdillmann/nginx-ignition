package formatter

import (
	"io"
	"text/template"

	"dillmann.com.br/nginx-ignition/tools/i18n/reader"
)

type TemplateBuilder func(file *reader.PropertiesFile) (string, any)

type templateFormatter struct {
	template string
}

func FromTemplate(template string) Formatter {
	return &templateFormatter{
		template: template,
	}
}

func (f *templateFormatter) Format(file *reader.PropertiesFile, writer io.Writer) error {
	tpl, err := template.New("").Parse(f.template)
	if err != nil {
		return err
	}

	return tpl.Execute(writer, file)
}
