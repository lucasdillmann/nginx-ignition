package host

import (
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"dillmann.com.br/nginx-ignition/core/host"
)

func newHostRequestDTO() hostRequestDTO {
	return hostRequestDTO{
		Enabled:           new(true),
		DefaultServer:     new(false),
		UseGlobalBindings: new(true),
		DomainNames:       []string{"example.com"},
		FeatureSet: &featureSetDTO{
			WebsocketsSupport:   new(true),
			HTTP2Support:        new(true),
			RedirectHTTPToHTTPS: new(true),
			StatsEnabled:        new(true),
		},
		Routes: []routeDTO{
			{
				Priority:   new(1),
				Enabled:    new(true),
				Type:       new(host.ProxyRouteType),
				SourcePath: new("/"),
				TargetURI:  new("http://backend"),
				Settings: &routeSettingsDTO{
					IncludeForwardHeaders:  new(true),
					ProxySslServerName:     new(true),
					KeepOriginalDomainName: new(true),
					IndexFile:              new("index.html"),
				},
			},
		},
	}
}

func newHost() *host.Host {
	return &host.Host{
		ID:                uuid.New(),
		Enabled:           true,
		DefaultServer:     false,
		UseGlobalBindings: true,
		DomainNames:       []string{"example.com"},
		FeatureSet: host.FeatureSet{
			WebsocketSupport:    true,
			HTTP2Support:        true,
			RedirectHTTPToHTTPS: true,
			StatsEnabled:        true,
		},
		Routes: []host.Route{
			{
				ID:         uuid.New(),
				Priority:   1,
				Enabled:    true,
				Type:       host.ProxyRouteType,
				SourcePath: "/",
				TargetURI:  new("http://backend"),
				Settings: host.RouteSettings{
					IncludeForwardHeaders:  true,
					ProxySSLServerName:     true,
					KeepOriginalDomainName: true,
					IndexFile:              new("index.html"),
				},
			},
		},
	}
}

func newHostPage() *pagination.Page[host.Host] {
	return pagination.Of([]host.Host{
		*newHost(),
	})
}
