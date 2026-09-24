package host

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/binding"
	host2 "github.com/lucasdillmann/nginx-ignition/internal/business/domain/host"
)

func cleanup(ctx context.Context, t *testing.T, repo host2.Repository) {
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

func newHost() *host2.Host {
	return &host2.Host{
		ID:                uuid.New(),
		DomainNames:       []string{"example.com"},
		Enabled:           true,
		DefaultServer:     false,
		UseGlobalBindings: true,
		Routes: []host2.Route{
			{
				ID:         uuid.New(),
				Priority:   10,
				Type:       host2.StaticResponseRouteType,
				SourcePath: "/",
				Settings: host2.RouteSettings{
					IncludeForwardHeaders:   true,
					IgnoreSSLErrors:         true,
					ProxySSLServerName:      false,
					KeepOriginalDomainName:  true,
					DirectoryListingEnabled: false,
					IndexFile:               new("index.html"),
					Custom:                  new("# Custom config"),
				},
				Response: &host2.RouteStaticResponse{
					StatusCode: 200,
					Payload:    new("OK"),
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
		FeatureSet: host2.FeatureSet{
			WebsocketSupport:    true,
			HTTP2Support:        true,
			RedirectHTTPToHTTPS: false,
			StatsEnabled:        true,
		},
		VPNs: []host2.VPN{},
	}
}
