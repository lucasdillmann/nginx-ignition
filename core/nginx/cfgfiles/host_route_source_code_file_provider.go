package cfgfiles

import (
	"dillmann.com.br/nginx-ignition/core/host"
	"fmt"
)

type hostRouteSourceCodeFileProvider struct{}

func newHostRouteSourceCodeFileProvider() *hostRouteSourceCodeFileProvider {
	return &hostRouteSourceCodeFileProvider{}
}

func (p *hostRouteSourceCodeFileProvider) provide(ctx *providerContext) ([]output, error) {
	var outputs []output
	for _, h := range ctx.hosts {
		if h.Enabled {
			outputs = append(outputs, p.buildSourceCodeFiles(h)...)
		}
	}

	return outputs, nil
}

func (p *hostRouteSourceCodeFileProvider) buildSourceCodeFiles(h *host.Host) []output {
	var outputs []output
	for _, r := range h.Routes {
		if r.Enabled && r.Type == host.SourceCodeRouteType && r.SourceCode.Language == host.JavascriptCodeLanguage {
			outputs = append(outputs, output{
				name:     fmt.Sprintf("host-%s-route-%d.js", h.ID, r.Priority),
				contents: r.SourceCode.Contents,
			})
		}
	}

	return outputs
}
