package provider

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/coreerror"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/i18n"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/host"
)

func Test_hostRouteSourceCodeProvider(t *testing.T) {
	t.Run("Provide", func(t *testing.T) {
		provider := &hostRouteSourceCodeProvider{}
		hostID := uuid.New()
		ctx := newProviderContext(t)
		ctx.SupportedFeatures.RunCodeType = DynamicSupportType
		ctx.Hosts = []host.Host{
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
		}

		files, err := provider.Provide(ctx)
		assert.NoError(t, err)
		assert.Len(t, files, 1)
		assert.Equal(t, fmt.Sprintf("host-%s-route-10.js", hostID), files[0].Name)
	})

	t.Run("BuildSourceCodeFiles", func(t *testing.T) {
		provider := &hostRouteSourceCodeProvider{}
		hostID := uuid.New()

		t.Run("generates javascript files when supported", func(t *testing.T) {
			ctx := newProviderContext(t)
			ctx.SupportedFeatures.RunCodeType = DynamicSupportType
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

			files, err := provider.buildSourceCodeFiles(ctx, h)
			assert.NoError(t, err)
			assert.Len(t, files, 1)
			assert.Equal(t, fmt.Sprintf("host-%s-route-10.js", hostID), files[0].Name)
			assert.Equal(t, "console.log('hi');", files[0].Contents)
		})

		t.Run("returns error when code execution is not supported", func(t *testing.T) {
			ctx := newProviderContext(t)
			ctx.SupportedFeatures.RunCodeType = NoneSupportType
			h := &host.Host{
				Routes: []host.Route{
					{
						Enabled: true,
						Type:    host.ExecuteCodeRouteType,
					},
				},
			}

			_, err := provider.buildSourceCodeFiles(ctx, h)
			assert.Error(t, err)
			var coreErr *coreerror.CoreError
			assert.ErrorAs(t, err, &coreErr)
			assert.Equal(t, i18n.K.CoreNginxCfgfilesHostRouteCodeNotEnabled, coreErr.Message.Key)
		})
	})
}
