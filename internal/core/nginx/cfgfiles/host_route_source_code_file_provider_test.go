package cfgfiles

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/lucasdillmann/nginx-ignition/internal/core/common/coreerror"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/i18n"
	"github.com/lucasdillmann/nginx-ignition/internal/core/host"
)

func Test_hostRouteSourceCodeFileProvider(t *testing.T) {
	t.Run("Provide", func(t *testing.T) {
		provider := &hostRouteSourceCodeFileProvider{}
		hostID := uuid.New()
		ctx := newProviderContext(t)
		ctx.supportedFeatures.RunCodeType = DynamicSupportType
		ctx.hosts = []host.Host{
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

		files, err := provider.provide(ctx)
		assert.NoError(t, err)
		assert.Len(t, files, 1)
		assert.Equal(t, fmt.Sprintf("host-%s-route-10.js", hostID), files[0].Name)
	})

	t.Run("BuildSourceCodeFiles", func(t *testing.T) {
		provider := &hostRouteSourceCodeFileProvider{}
		hostID := uuid.New()

		t.Run("generates javascript files when supported", func(t *testing.T) {
			ctx := newProviderContext(t)
			ctx.supportedFeatures.RunCodeType = DynamicSupportType
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
			ctx.supportedFeatures.RunCodeType = NoneSupportType
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
