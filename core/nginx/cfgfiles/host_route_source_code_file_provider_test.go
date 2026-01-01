package cfgfiles

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"dillmann.com.br/nginx-ignition/core/host"
)

func TestHostRouteSourceCodeFileProvider_Provide(t *testing.T) {
	p := &hostRouteSourceCodeFileProvider{}
	hostID := uuid.New()
	ctx := &providerContext{
		supportedFeatures: &SupportedFeatures{
			RunCodeType: DynamicSupportType,
		},
		hosts: []host.Host{
			{
				ID: hostID,
				Routes: []host.Route{
					{
						Enabled:  true,
						Priority: 10,
						Type:     host.ExecuteCodeRouteType,
						SourceCode: &host.RouteSourceCode{
							Language: host.JavascriptCodeLanguage,
							Contents: "console.log('hi');",
						},
					},
				},
			},
		},
	}

	files, err := p.provide(ctx)
	assert.NoError(t, err)
	assert.Len(t, files, 1)
	assert.Equal(t, fmt.Sprintf("host-%s-route-10.js", hostID), files[0].Name)
}

func TestHostRouteSourceCodeFileProvider_BuildSourceCodeFiles(t *testing.T) {
	p := &hostRouteSourceCodeFileProvider{}
	hostID := uuid.New()

	t.Run("generates javascript files when supported", func(t *testing.T) {
		ctx := &providerContext{
			supportedFeatures: &SupportedFeatures{
				RunCodeType: DynamicSupportType,
			},
		}
		h := &host.Host{
			ID: hostID,
			Routes: []host.Route{
				{
					Enabled:  true,
					Priority: 10,
					Type:     host.ExecuteCodeRouteType,
					SourceCode: &host.RouteSourceCode{
						Language: host.JavascriptCodeLanguage,
						Contents: "console.log('hi');",
					},
				},
			},
		}

		files, err := p.buildSourceCodeFiles(ctx, h)
		assert.NoError(t, err)
		assert.Len(t, files, 1)
		assert.Equal(t, fmt.Sprintf("host-%s-route-10.js", hostID), files[0].Name)
		assert.Equal(t, "console.log('hi');", files[0].Contents)
	})

	t.Run("returns error when code execution is not supported", func(t *testing.T) {
		ctx := &providerContext{
			supportedFeatures: &SupportedFeatures{
				RunCodeType: NoneSupportType,
			},
		}
		h := &host.Host{
			Routes: []host.Route{
				{
					Enabled: true,
					Type:    host.ExecuteCodeRouteType,
				},
			},
		}

		_, err := p.buildSourceCodeFiles(ctx, h)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not enabled")
	})
}
