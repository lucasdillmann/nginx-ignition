package host

import (
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/binding"
	"dillmann.com.br/nginx-ignition/core/common/ptr"
)

func newHost() *Host {
	return &Host{
		ID:          uuid.New(),
		Enabled:     true,
		DomainNames: []string{"example.com"},
		Bindings: []binding.Binding{
			{
				Type: binding.HTTPBindingType,
				IP:   "0.0.0.0",
				Port: 80,
			},
		},
		Routes: []Route{
			{
				ID:         uuid.New(),
				Enabled:    true,
				Priority:   0,
				SourcePath: "/",
				Type:       StaticResponseRouteType,
				Response: &RouteStaticResponse{
					StatusCode: 200,
					Payload:    ptr.Of("OK"),
				},
			},
		},
		FeatureSet: FeatureSet{
			StatsEnabled: true,
		},
	}
}
