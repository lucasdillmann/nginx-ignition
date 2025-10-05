package cfgfiles

import (
	"fmt"

	"dillmann.com.br/nginx-ignition/core/common/core_error"
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
			files, err := p.buildSourceCodeFiles(ctx, h)
			if err != nil {
				return nil, err
			}

			outputs = append(outputs, *files...)
		}
	}

	return outputs, nil
}

func (p *hostRouteSourceCodeFileProvider) buildSourceCodeFiles(
	ctx *providerContext,
	h *host.Host,
) (*[]File, error) {
	var outputs []File

	for _, r := range h.Routes {
		if !r.Enabled || r.Type != host.ExecuteCodeRouteType {
			continue
		}

		if !ctx.supportedFeatures.RunCode {
			return nil, core_error.New(
				"Unable to generate the host route source code files: Run code support is not enabled in the nginx server.",
				false,
			)
		}

		if r.SourceCode.Language == host.JavascriptCodeLanguage {
			outputs = append(outputs, File{
				Name:     fmt.Sprintf("host-%s-route-%d.js", h.ID, r.Priority),
				Contents: r.SourceCode.Contents,
			})
		}
	}

	return &outputs, nil
}
