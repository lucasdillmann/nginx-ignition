package cfgfiles

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"dillmann.com.br/nginx-ignition/core/common/ptr"
	"dillmann.com.br/nginx-ignition/core/host"
)

func Test_HostRouteStaticResponseFileProvider_Provide(t *testing.T) {
	p := &hostRouteStaticResponseFileProvider{}
	hostID := uuid.New()
	ctx := &providerContext{
		hosts: []host.Host{
			{
				ID: hostID,
				Routes: []host.Route{
					{
						Enabled:  true,
						Priority: 10,
						Type:     host.StaticResponseRouteType,
						Response: &host.RouteStaticResponse{
							Payload: ptr.Of("hello world"),
						},
					},
				},
			},
		},
	}

	files, err := p.provide(ctx)
	assert.NoError(t, err)
	assert.Len(t, files, 1)
	assert.Equal(t, fmt.Sprintf("host-%s-route-10.payload", hostID), files[0].Name)
}

func Test_HostRouteStaticResponseFileProvider_BuildStaticResponseFiles(t *testing.T) {
	p := &hostRouteStaticResponseFileProvider{}
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
						Payload: ptr.Of("hello world"),
					},
				},
				{
					Enabled:  false,
					Priority: 20,
					Type:     host.StaticResponseRouteType,
				},
			},
		}

		files := p.buildStaticResponseFiles(h)
		assert.Len(t, files, 1)
		assert.Equal(t, fmt.Sprintf("host-%s-route-10.payload", hostID), files[0].Name)
		assert.Equal(t, "hello world", files[0].Contents)
	})
}
