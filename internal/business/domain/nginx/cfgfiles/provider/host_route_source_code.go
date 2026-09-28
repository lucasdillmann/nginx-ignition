package provider

import (
	"fmt"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/coreerror"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/i18n"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/host"
)

type hostRouteSourceCodeProvider struct{}

func newHostRouteSourceCodeProvider() *hostRouteSourceCodeProvider {
	return &hostRouteSourceCodeProvider{}
}

func (p *hostRouteSourceCodeProvider) Provide(ctx *Context) ([]File, error) {
	outputs := make([]File, 0)

	for _, h := range ctx.Hosts {
		files, err := p.buildSourceCodeFiles(ctx, &h)
		if err != nil {
			return nil, err
		}

		outputs = append(outputs, files...)
	}

	return outputs, nil
}

func (p *hostRouteSourceCodeProvider) buildSourceCodeFiles(
	ctx *Context,
	h *host.Host,
) ([]File, error) {
	outputs := make([]File, 0)

	for _, r := range h.Routes {
		if !r.Enabled || r.Type != host.ExecuteCodeRouteType {
			continue
		}

		if ctx.SupportedFeatures.RunCodeType == NoneSupportType {
			return nil, coreerror.New(
				i18n.M(ctx.Context, i18n.K.CoreNginxCfgfilesHostRouteCodeNotEnabled),
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
