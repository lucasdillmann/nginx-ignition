package host

import (
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"dillmann.com.br/nginx-ignition/core/common/ptr"
	"dillmann.com.br/nginx-ignition/core/host"
)

func newHostRequestDTO() hostRequestDTO {
	return hostRequestDTO{
		Enabled:           ptr.Of(true),
		DefaultServer:     ptr.Of(false),
		UseGlobalBindings: ptr.Of(true),
		DomainNames:       []string{"example.com"},
		FeatureSet: &featureSetDTO{
			WebsocketsSupport:   ptr.Of(true),
			HTTP2Support:        ptr.Of(true),
			RedirectHTTPToHTTPS: ptr.Of(true),
			StatsEnabled:        ptr.Of(true),
		},
		Routes: []routeDTO{
			{
				Priority:   ptr.Of(1),
				Enabled:    ptr.Of(true),
				Type:       ptr.Of(host.ProxyRouteType),
				SourcePath: ptr.Of("/"),
				TargetURI:  ptr.Of("http://backend"),
				Settings: &routeSettingsDTO{
					IncludeForwardHeaders:  ptr.Of(true),
					ProxySslServerName:     ptr.Of(true),
					KeepOriginalDomainName: ptr.Of(true),
					IndexFile:              ptr.Of("index.html"),
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
				TargetURI:  ptr.Of("http://backend"),
				Settings: host.RouteSettings{
					IncludeForwardHeaders:  true,
					ProxySSLServerName:     true,
					KeepOriginalDomainName: true,
					IndexFile:              ptr.Of("index.html"),
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
