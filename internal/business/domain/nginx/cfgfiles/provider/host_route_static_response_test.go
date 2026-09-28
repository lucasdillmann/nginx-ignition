package provider

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/host"
)

func Test_hostRouteStaticResponseProvider(t *testing.T) {
	t.Run("Provide", func(t *testing.T) {
		provider := &hostRouteStaticResponseProvider{}
		hostID := uuid.New()
		ctx := newProviderContext(t)
		ctx.Hosts = []host.Host{
			{
				ID: hostID,
				Routes: []host.Route{
					{
						Enabled:  true,
						Priority: 10,
						Type:     host.StaticResponseRouteType,
						Response: &host.RouteStaticResponse{
							Payload: new("hello world"),
						},
					},
				},
			},
		}

		files, err := provider.Provide(ctx)
		assert.NoError(t, err)
		assert.Len(t, files, 1)
		assert.Equal(t, fmt.Sprintf("host-%s-route-10.payload", hostID), files[0].Name)
	})

	t.Run("BuildStaticResponseFiles", func(t *testing.T) {
		provider := &hostRouteStaticResponseProvider{}
		hostID := uuid.New()

		t.Run("generates files for enabled static routes", func(t *testing.T) {
			h := &host.Host{
				ID: hostID,
				Routes: []host.Route{
					{
						Enabled:  true,
						Priority: 10,
						Type:     host.StaticResponseRouteType,
						Response: &host.RouteStaticResponse{
							Payload: new("hello world"),
						},
					},
					{
						Enabled:  false,
						Priority: 20,
						Type:     host.StaticResponseRouteType,
					},
				},
			}

			files := provider.buildStaticResponseFiles(h)
			assert.Len(t, files, 1)
			assert.Equal(t, fmt.Sprintf("host-%s-route-10.payload", hostID), files[0].Name)
			assert.Equal(t, "hello world", files[0].Contents)
		})
	})
}
