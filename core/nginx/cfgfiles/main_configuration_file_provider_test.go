package cfgfiles

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"dillmann.com.br/nginx-ignition/core/cache"
	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"dillmann.com.br/nginx-ignition/core/common/ptr"
	"dillmann.com.br/nginx-ignition/core/host"
	"dillmann.com.br/nginx-ignition/core/settings"
	"dillmann.com.br/nginx-ignition/core/stream"
)

func Test_mainConfigurationFileProvider(t *testing.T) {
	t.Run("Provide", func(t *testing.T) {
		provider := &mainConfigurationFileProvider{
			config: configuration.New(),
		}

		paths := newPaths()
		mockSettings := newSettings()

		t.Run("successfully generates basic config", func(t *testing.T) {
			ctx := newProviderContext(t)
			ctx.paths = paths
			ctx.cfg = mockSettings
			ctx.supportedFeatures.StreamType = NoneSupportType
			ctx.supportedFeatures.RunCodeType = NoneSupportType

			files, err := provider.provide(ctx)
			assert.NoError(t, err)
			assert.Len(t, files, 1)
			assert.Contains(t, files[0].Contents, "worker_processes 1;")
			assert.NotContains(t, files[0].Contents, "load_module")
			assert.NotContains(t, files[0].Contents, "stream {")
		})

		t.Run("includes dynamic modules and stream block when enabled", func(t *testing.T) {
			ctx := newProviderContext(t)
			ctx.paths = paths
			ctx.cfg = mockSettings
			ctx.supportedFeatures.StreamType = DynamicSupportType
			ctx.supportedFeatures.RunCodeType = DynamicSupportType
			ctx.streams = []stream.Stream{{ID: uuid.New()}}

			files, err := provider.provide(ctx)
			assert.NoError(t, err)
			assert.Contains(t, files[0].Contents, "load_module modules/ngx_stream_module.so;")
			assert.Contains(t, files[0].Contents, "load_module modules/ngx_http_js_module.so;")
			assert.Contains(t, files[0].Contents, "stream {")
			assert.Contains(t, files[0].Contents, "include \"/etc/nginx/stream-")
		})

		t.Run("includes custom configuration", func(t *testing.T) {
			mockSettings.Nginx.Custom = ptr.Of("custom_directive on;")
			ctx := newProviderContext(t)
			ctx.paths = paths
			ctx.cfg = mockSettings
			ctx.supportedFeatures.StreamType = NoneSupportType
			ctx.supportedFeatures.RunCodeType = NoneSupportType

			files, err := provider.provide(ctx)
			assert.NoError(t, err)
			assert.Contains(t, files[0].Contents, "custom_directive on;")
		})
	})

	t.Run("getErrorLogPath", func(t *testing.T) {
		provider := &mainConfigurationFileProvider{
			config: configuration.New(),
		}
		paths := &Paths{
			Logs: "/var/log/nginx/",
		}

		t.Run("returns off when disabled", func(t *testing.T) {
			logs := &settings.NginxLogsSettings{
				ServerLogsEnabled: false,
			}
			assert.Equal(t, "off", provider.getErrorLogPath(paths, logs))
		})

		t.Run("returns path and level when enabled", func(t *testing.T) {
			logs := &settings.NginxLogsSettings{
				ServerLogsEnabled: true,
				ServerLogsLevel:   settings.WarnLogLevel,
			}
			assert.Equal(
				t,
				"\"/var/log/nginx/main.log\" warn",
				provider.getErrorLogPath(paths, logs),
			)
		})
	})

	t.Run("getHostIncludes", func(t *testing.T) {
		provider := &mainConfigurationFileProvider{
			config: configuration.New(),
		}
		paths := &Paths{
			Config: "/etc/nginx/",
		}
		id1 := uuid.New()
		id2 := uuid.New()

		t.Run("returns empty string for no hosts", func(t *testing.T) {
			assert.Equal(t, "", provider.getHostIncludes(paths, nil))
		})

		t.Run("returns include directives for multiple hosts", func(t *testing.T) {
			hosts := []host.Host{
				{
					ID: id1,
				},
				{
					ID: id2,
				},
			}
			result := provider.getHostIncludes(paths, hosts)
			assert.Contains(t, result, fmt.Sprintf("include \"/etc/nginx/host-%s.conf\";", id1))
			assert.Contains(t, result, fmt.Sprintf("include \"/etc/nginx/host-%s.conf\";", id2))
		})
	})

	t.Run("getStreamIncludes", func(t *testing.T) {
		provider := &mainConfigurationFileProvider{
			config: configuration.New(),
		}
		paths := &Paths{
			Config: "/etc/nginx/",
		}
		id1 := uuid.New()

		t.Run("returns include directives for streams", func(t *testing.T) {
			streams := []stream.Stream{
				{
					ID: id1,
				},
			}
			result := provider.getStreamIncludes(paths, streams)
			assert.Equal(t, fmt.Sprintf("include \"/etc/nginx/stream-%s.conf\";", id1), result)
		})
	})

	t.Run("getCacheDefinitions", func(t *testing.T) {
		provider := &mainConfigurationFileProvider{
			config: configuration.New(),
		}
		paths := &Paths{
			Cache: "/var/cache/nginx/",
		}
		id1 := uuid.New()

		t.Run("generates proxy_cache_path with various options", func(t *testing.T) {
			caches := []cache.Cache{
				{
					ID:              id1,
					InactiveSeconds: ptr.Of(3600),
					MaximumSizeMB:   ptr.Of(1024),
				},
			}
			result := provider.getCacheDefinitions(paths, caches)
			assert.Contains(t, result, "proxy_cache_path \"/var/cache/nginx/")
			assert.Contains(t, result, "inactive=3600s")
			assert.Contains(t, result, "max_size=1024m")
			assert.Contains(t, result, "keys_zone=cache_")
		})

		t.Run("uses custom storage path if provided", func(t *testing.T) {
			customPath := "/mnt/ssd/cache"
			caches := []cache.Cache{
				{
					ID:          id1,
					StoragePath: &customPath,
				},
			}
			result := provider.getCacheDefinitions(paths, caches)
			assert.Contains(t, result, "proxy_cache_path \"/mnt/ssd/cache\"")
		})

		t.Run("generates basic config when optional fields are nil", func(t *testing.T) {
			caches := []cache.Cache{
				{
					ID: id1,
				},
			}
			result := provider.getCacheDefinitions(paths, caches)
			assert.Contains(t, result, "proxy_cache_path \"/var/cache/nginx/")
			assert.NotContains(t, result, "inactive=")
			assert.NotContains(t, result, "max_size=")
		})
	})

	t.Run("getStatsDefinitions", func(t *testing.T) {
		provider := &mainConfigurationFileProvider{
			config: configuration.New(),
		}
		paths := &Paths{
			Base: "/etc/nginx/",
		}

		t.Run("returns empty string when disabled", func(t *testing.T) {
			cfg := &settings.NginxStatsSettings{
				Enabled: false,
			}
			result, err := provider.getStatsDefinitions(paths, cfg)
			assert.NoError(t, err)
			assert.Equal(t, "", result)
		})

		t.Run("returns empty string when nil", func(t *testing.T) {
			result, err := provider.getStatsDefinitions(paths, nil)
			assert.NoError(t, err)
			assert.Equal(t, "", result)
		})

		t.Run("generates base config when enabled", func(t *testing.T) {
			cfg := &settings.NginxStatsSettings{
				Enabled:       true,
				MaximumSizeMB: 10,
				Persistent:    false,
			}
			result, err := provider.getStatsDefinitions(paths, cfg)
			assert.NoError(t, err)
			assert.Contains(
				t,
				result,
				"vhost_traffic_status_zone shared:nginx-ignition-traffic-stats:10m;",
			)
			assert.Contains(t, result, "vhost_traffic_status_filter_by_host on;")
			assert.Contains(t, result, "vhost_traffic_status_stats_by_upstream on;")
			assert.Contains(t, result, "server {")
			assert.Contains(t, result, "listen unix:/etc/nginx/traffic-stats.socket;")
			assert.NotContains(t, result, "vhost_traffic_status_dump")
		})

		t.Run("includes persistent dump with default path", func(t *testing.T) {
			cfg := &settings.NginxStatsSettings{
				Enabled:       true,
				MaximumSizeMB: 10,
				Persistent:    true,
			}
			result, err := provider.getStatsDefinitions(paths, cfg)
			assert.NoError(t, err)
			assert.Contains(
				t,
				result,
				"vhost_traffic_status_dump \"/tmp/nginx-ignition/data/stats.db\" 5s;",
			)
		})

		t.Run("includes persistent dump with custom path", func(t *testing.T) {
			cfg := &settings.NginxStatsSettings{
				Enabled:          true,
				MaximumSizeMB:    10,
				Persistent:       true,
				DatabaseLocation: ptr.Of("/var/lib/nginx/stats.db"),
			}
			result, err := provider.getStatsDefinitions(paths, cfg)
			assert.NoError(t, err)
			assert.Contains(t, result, "vhost_traffic_status_dump \"/var/lib/nginx/stats.db\" 5s;")
		})
	})
}
