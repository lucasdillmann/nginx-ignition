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
	indentLevel := 0
	indentValue := strings.Repeat(" ", 4)
	originalLines := strings.Split(f.Contents, "\n")
	formattedLines := make([]string, 0, len(originalLines))

	for index, line := range originalLines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" && index == 0 {
			continue
		}

		if strings.HasSuffix(trimmed, "}") && !strings.Contains(trimmed, "{") {
			if indentLevel > 0 {
				indentLevel--
			}
		}

		indentedLine := strings.Repeat(indentValue, indentLevel) + trimmed
		formattedLines = append(formattedLines, indentedLine)

		if strings.HasSuffix(trimmed, "{") && !strings.HasSuffix(trimmed, "};") {
			indentLevel++
		}
	}

	return strings.Join(formattedLines, "\n")
}
