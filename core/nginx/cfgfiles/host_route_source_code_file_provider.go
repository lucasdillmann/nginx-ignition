package cfgfiles

import (
	"fmt"

	"dillmann.com.br/nginx-ignition/core/host"
)

type hostRouteSourceCodeFileProvider struct{}

func newHostRouteSourceCodeFileProvider() *hostRouteSourceCodeFileProvider {
	return &hostRouteSourceCodeFileProvider{}
}

func (p *hostRouteSourceCodeFileProvider) provide(ctx *providerContext) ([]File, error) {
	var outputs []File
	for _, h := range ctx.hosts {
		if h.Enabled {
			outputs = append(outputs, p.buildSourceCodeFiles(h)...)
		}
	}

	return outputs, nil
}

func (p *hostRouteSourceCodeFileProvider) buildSourceCodeFiles(h *host.Host) []File {
	var outputs []File
	for _, r := range h.Routes {
		if r.Enabled && r.Type == host.ExecuteCodeRouteType && r.SourceCode.Language == host.JavascriptCodeLanguage {
			outputs = append(outputs, File{
				Name:     fmt.Sprintf("host-%s-route-%d.js", h.ID, r.Priority),
				Contents: r.SourceCode.Contents,
			})
		}
	}

	return outputs
}
