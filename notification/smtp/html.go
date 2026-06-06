package smtp

import (
	"bytes"
	"strings"
	"text/template"

	"dillmann.com.br/nginx-ignition/core/notification"
)

const htmlBodyTemplate = `<html><body>
<p>{{html .Summary}}</p>
{{- range .Sections}}
{{- if .Title}}<h3>{{html .Title}}</h3>{{end}}
<p>{{html .Body}}</p>
{{- end}}
{{- if .Actions}}
<ul>
{{- range .Actions}}
<li><a href="{{html .URL}}">{{html .Label}}</a></li>
{{- end}}
</ul>
{{- end}}
</body></html>`

var htmlBodyTmpl = template.Must(
	template.New("htmlBody").Funcs(template.FuncMap{
		"html": htmlEscape,
	}).Parse(htmlBodyTemplate),
)

type htmlBodyData struct {
	Summary  string
	Sections []htmlSectionData
	Actions  []htmlActionData
}

type htmlSectionData struct {
	Title *string
	Body  string
}

type htmlActionData struct {
	Label string
	URL   string
}

func formatHTMLBody(deliverable notification.Deliverable) string {
	data := htmlBodyData{
		Summary:  deliverable.Summary,
		Sections: make([]htmlSectionData, len(deliverable.Sections)),
		Actions:  make([]htmlActionData, len(deliverable.Actions)),
	}

	for index, section := range deliverable.Sections {
		data.Sections[index] = htmlSectionData{
			Title: section.Title,
			Body:  section.Body,
		}
	}

	for index, action := range deliverable.Actions {
		data.Actions[index] = htmlActionData{
			Label: action.Label,
			URL:   action.URL,
		}
	}

	var buffer bytes.Buffer
	_ = htmlBodyTmpl.Execute(&buffer, data)
	return buffer.String()
}

func htmlEscape(value string) string {
	replacer := strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
		"\"", "&quot;",
	)
	return replacer.Replace(value)
}
