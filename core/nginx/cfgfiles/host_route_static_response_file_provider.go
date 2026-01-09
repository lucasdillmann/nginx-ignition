package cfgfiles

import (
	"fmt"

	"dillmann.com.br/nginx-ignition/core/host"
)

type hostRouteStaticResponseFileProvider struct{}

func newHostRouteStaticResponseFileProvider() *hostRouteStaticResponseFileProvider {
	return &hostRouteStaticResponseFileProvider{}
}

func (p *hostRouteStaticResponseFileProvider) provide(ctx *providerContext) ([]File, error) {
	outputs := make([]File, 0, len(ctx.hosts))

	for _, h := range ctx.hosts {
		outputs = append(outputs, p.buildStaticResponseFiles(&h)...)
	}

	return outputs, nil
}

func (p *hostRouteStaticResponseFileProvider) buildStaticResponseFiles(h *host.Host) []File {
	outputs := make([]File, 0)

	for _, r := range h.Routes {
		if !r.Enabled || r.Type != host.StaticResponseRouteType {
			continue
		}

		var contents string
		if r.Response.Payload != nil {
			contents = *r.Response.Payload
		}

		outputs = append(outputs, File{
			Name:     fmt.Sprintf("host-%s-route-%d.payload", h.ID, r.Priority),
			Contents: contents,
		})
	}

	return outputs
}
