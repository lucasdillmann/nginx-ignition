package cfgfiles

import (
	"context"
	"strings"

	"dillmann.com.br/nginx-ignition/core/host"
	"dillmann.com.br/nginx-ignition/core/stream"
)

type providerContext struct {
	context context.Context
	paths   *Paths
	hosts   []*host.Host
	streams []*stream.Stream
}

type Paths struct {
	Config string
	Logs   string
}

type fileProvider interface {
	provide(ctx *providerContext) ([]File, error)
}

type File struct {
	Name     string
	Contents string
}

func (f *File) FormattedContents() string {
	if f.Contents == "" {
		return ""
	}

	lines := strings.Split(f.Contents, "\n")

	var formatted []string
	indentLevel := 0
	indentValue := strings.Repeat(" ", 4)

	for index, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			if index != 0 {
				formatted = append(formatted, "")
			}

			continue
		}

		if strings.HasSuffix(trimmed, "}") && !strings.Contains(trimmed, "{") {
			if indentLevel > 0 {
				indentLevel--
			}
		}

		indentedLine := strings.Repeat(indentValue, indentLevel) + trimmed
		formatted = append(formatted, indentedLine)

		if strings.HasSuffix(trimmed, "{") && !strings.HasSuffix(trimmed, "};") {
			indentLevel++
		}
	}

	return strings.Join(formatted, "\n")
}
