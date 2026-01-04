package host

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"dillmann.com.br/nginx-ignition/core/binding"
	"dillmann.com.br/nginx-ignition/core/common/ptr"
	"dillmann.com.br/nginx-ignition/core/host"
)

func cleanup(ctx context.Context, t *testing.T, repo host.Repository) {
	result, err := repo.(*repository).database.Unwrap().Query("SELECT id FROM host")
	require.NoError(t, err)

	defer result.Close()

	ids := make([]string, 0)
	for result.Next() {
		var id string

		err = result.Scan(&id)
		require.NoError(t, err)

		ids = append(ids, id)
	}

	result.Close()

	for _, id := range ids {
		err = repo.DeleteByID(ctx, uuid.MustParse(id))
		require.NoError(t, err)
	}
}

func newHost() *host.Host {
	return &host.Host{
		ID:                uuid.New(),
		DomainNames:       []string{"example.com"},
		Enabled:           true,
		DefaultServer:     false,
		UseGlobalBindings: true,
		Routes: []host.Route{
			{
				ID:         uuid.New(),
				Priority:   10,
				Type:       host.StaticResponseRouteType,
				SourcePath: "/",
				Settings: host.RouteSettings{
					IncludeForwardHeaders:   true,
					ProxySSLServerName:      false,
					KeepOriginalDomainName:  true,
					DirectoryListingEnabled: false,
					IndexFile:               ptr.Of("index.html"),
					Custom:                  ptr.Of("# Custom config"),
				},
				Response: &host.RouteStaticResponse{
					StatusCode: 200,
					Payload:    ptr.Of("OK"),
				},
			},
		},
		Bindings: []binding.Binding{
			{
				ID:   uuid.New(),
				Type: binding.HTTPBindingType,
				IP:   "0.0.0.0",
				Port: 8080,
			},
		},
		FeatureSet: host.FeatureSet{
			WebsocketSupport:    true,
			HTTP2Support:        true,
			RedirectHTTPToHTTPS: false,
		},
		VPNs: []host.VPN{},
	}
}
