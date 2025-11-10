package cfgfiles

import (
	"fmt"

	"dillmann.com.br/nginx-ignition/core/common/coreerror"
	"dillmann.com.br/nginx-ignition/core/host"
)

type hostRouteSourceCodeFileProvider struct{}

func newHostRouteSourceCodeFileProvider() *hostRouteSourceCodeFileProvider {
	return &hostRouteSourceCodeFileProvider{}
}

func (p *hostRouteSourceCodeFileProvider) provide(ctx *providerContext) ([]File, error) {
	var outputs []File

	for _, h := range ctx.hosts {
		files, err := p.buildSourceCodeFiles(ctx, h)
		if err != nil {
			return nil, err
		}

		outputs = append(outputs, *files...)
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

		if ctx.supportedFeatures.RunCodeType == NoneSupportType {
			return nil, coreerror.New(
				"Unable to generate the host route source code files: Support for JavaScript and/or Lua "+
					"code is not enabled in the nginx server and at least one code execution host route is enabled.",
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
