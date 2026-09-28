package provider

import (
	"fmt"

	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/host"
)

type hostRouteStaticResponseProvider struct{}

func newHostRouteStaticResponseProvider() *hostRouteStaticResponseProvider {
	return &hostRouteStaticResponseProvider{}
}

func (p *hostRouteStaticResponseProvider) Provide(ctx *Context) ([]File, error) {
	outputs := make([]File, 0, len(ctx.Hosts))

	for _, h := range ctx.Hosts {
		outputs = append(outputs, p.buildStaticResponseFiles(&h)...)
	}

	return outputs, nil
}

func (p *hostRouteStaticResponseProvider) buildStaticResponseFiles(h *host.Host) []File {
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
