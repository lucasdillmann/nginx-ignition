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
	var outputs []File

	for _, h := range ctx.hosts {
		files, err := p.buildStaticResponseFiles(&h)
		if err != nil {
			return nil, err
		}

		outputs = append(outputs, files...)
	}

	return outputs, nil
}

func (p *hostRouteStaticResponseFileProvider) buildStaticResponseFiles(h *host.Host) ([]File, error) {
	var outputs []File

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

	return outputs, nil
}
