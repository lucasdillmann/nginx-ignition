package cfgfiles

import (
	"fmt"

	"dillmann.com.br/nginx-ignition/core/common/coreerror"
	"dillmann.com.br/nginx-ignition/core/host"
	"dillmann.com.br/nginx-ignition/i18n"
)

type hostRouteSourceCodeFileProvider struct{}

func newHostRouteSourceCodeFileProvider() *hostRouteSourceCodeFileProvider {
	return &hostRouteSourceCodeFileProvider{}
}

func (p *hostRouteSourceCodeFileProvider) provide(ctx *providerContext) ([]File, error) {
	outputs := make([]File, 0)

	for _, h := range ctx.hosts {
		files, err := p.buildSourceCodeFiles(ctx, &h)
		if err != nil {
			return nil, err
		}

		outputs = append(outputs, files...)
	}

	return outputs, nil
}

func (p *hostRouteSourceCodeFileProvider) buildSourceCodeFiles(
	ctx *providerContext,
	h *host.Host,
) ([]File, error) {
	outputs := make([]File, 0)

	for _, r := range h.Routes {
		if !r.Enabled || r.Type != host.ExecuteCodeRouteType {
			continue
		}

		if ctx.supportedFeatures.RunCodeType == NoneSupportType {
			return nil, coreerror.New(
				i18n.M(ctx.context, i18n.K.NginxErrorHostRouteCodeNotEnabled),
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

	return outputs, nil
}
