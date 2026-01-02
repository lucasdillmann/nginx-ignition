package cfgfiles

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"dillmann.com.br/nginx-ignition/core/cache"
	"dillmann.com.br/nginx-ignition/core/common/ptr"
	"dillmann.com.br/nginx-ignition/core/host"
	"dillmann.com.br/nginx-ignition/core/settings"
	"dillmann.com.br/nginx-ignition/core/stream"
)

func Test_MainConfigurationFileProvider_Provide(t *testing.T) {
	p := &mainConfigurationFileProvider{}
	paths := &Paths{
		Base:   "/",
		Config: "/etc/nginx/",
		Logs:   "/var/log/nginx/",
		Cache:  "/var/cache/nginx/",
	}

	mockSettings := &settings.Settings{
		Nginx: &settings.NginxSettings{
			RuntimeUser:       "nginx",
			WorkerProcesses:   1,
			WorkerConnections: 1024,
			Timeouts: &settings.NginxTimeoutsSettings{
				Keepalive:  65,
				Connect:    60,
				Read:       60,
				Send:       60,
				ClientBody: 60,
			},
			Buffers: &settings.NginxBuffersSettings{
				ClientBodyKb:   16,
				ClientHeaderKb: 1,
				LargeClientHeader: &settings.NginxBufferSize{
					Amount: 4,
					SizeKb: 8,
				},
				Output: &settings.NginxBufferSize{
					Amount: 2,
					SizeKb: 32,
				},
			},
			Logs: &settings.NginxLogsSettings{
				ServerLogsEnabled: true,
				ServerLogsLevel:   settings.WarnLogLevel,
			},
		},
	}

	p.settingsCommands = &settings.Commands{
		Get: func(_ context.Context) (*settings.Settings, error) {
			return mockSettings, nil
		},
	}

	t.Run("successfully generates basic config", func(t *testing.T) {
		ctx := &providerContext{
			context: context.Background(),
			paths:   paths,
			supportedFeatures: &SupportedFeatures{
				StreamType:  NoneSupportType,
				RunCodeType: NoneSupportType,
			},
		}

		files, err := p.provide(ctx)
		assert.NoError(t, err)
		assert.Len(t, files, 1)
		assert.Contains(t, files[0].Contents, "worker_processes 1;")
		assert.NotContains(t, files[0].Contents, "load_module")
		assert.NotContains(t, files[0].Contents, "stream {")
	})

	t.Run("includes dynamic modules and stream block when enabled", func(t *testing.T) {
		ctx := &providerContext{
			context: context.Background(),
			paths:   paths,
			supportedFeatures: &SupportedFeatures{
				StreamType:  DynamicSupportType,
				RunCodeType: DynamicSupportType,
			},
			streams: []stream.Stream{{ID: uuid.New()}},
		}

		files, err := p.provide(ctx)
		assert.NoError(t, err)
		assert.Contains(t, files[0].Contents, "load_module modules/ngx_stream_module.so;")
		assert.Contains(t, files[0].Contents, "load_module modules/ngx_http_js_module.so;")
		assert.Contains(t, files[0].Contents, "stream {")
		assert.Contains(t, files[0].Contents, "include /etc/nginx/stream-")
	})

	t.Run("includes custom configuration", func(t *testing.T) {
		mockSettings.Nginx.Custom = ptr.Of("custom_directive on;")
		ctx := &providerContext{
			context: context.Background(),
			paths:   paths,
			supportedFeatures: &SupportedFeatures{
				StreamType:  NoneSupportType,
				RunCodeType: NoneSupportType,
			},
		}

		files, err := p.provide(ctx)
		assert.NoError(t, err)
		assert.Contains(t, files[0].Contents, "custom_directive on;")
	})

	t.Run("returns error when settingsCommands fails", func(t *testing.T) {
		p.settingsCommands = &settings.Commands{
			Get: func(_ context.Context) (*settings.Settings, error) {
				return nil, assert.AnError
			},
		}
		ctx := &providerContext{context: context.Background()}
		_, err := p.provide(ctx)
		assert.ErrorIs(t, err, assert.AnError)
	})
}

func Test_MainConfigurationFileProvider_GetErrorLogPath(t *testing.T) {
	p := &mainConfigurationFileProvider{}
	paths := &Paths{
		Logs: "/var/log/nginx/",
	}

	t.Run("returns off when disabled", func(t *testing.T) {
		logs := &settings.NginxLogsSettings{
			ServerLogsEnabled: false,
		}
		assert.Equal(t, "off", p.getErrorLogPath(paths, logs))
	})

	t.Run("returns path and level when enabled", func(t *testing.T) {
		logs := &settings.NginxLogsSettings{
			ServerLogsEnabled: true,
			ServerLogsLevel:   settings.WarnLogLevel,
		}
		assert.Equal(t, "/var/log/nginx/main.log warn", p.getErrorLogPath(paths, logs))
	})
}

func Test_MainConfigurationFileProvider_GetHostIncludes(t *testing.T) {
	p := &mainConfigurationFileProvider{}
	paths := &Paths{
		Config: "/etc/nginx/",
	}
	id1 := uuid.New()
	id2 := uuid.New()

	t.Run("returns empty string for no hosts", func(t *testing.T) {
		assert.Equal(t, "", p.getHostIncludes(paths, nil))
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
		result := p.getHostIncludes(paths, hosts)
		assert.Contains(t, result, fmt.Sprintf("include /etc/nginx/host-%s.conf;", id1))
		assert.Contains(t, result, fmt.Sprintf("include /etc/nginx/host-%s.conf;", id2))
	})
}

func Test_MainConfigurationFileProvider_GetStreamIncludes(t *testing.T) {
	p := &mainConfigurationFileProvider{}
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
		result := p.getStreamIncludes(paths, streams)
		assert.Equal(t, fmt.Sprintf("include /etc/nginx/stream-%s.conf;", id1), result)
	})
}

func Test_MainConfigurationFileProvider_GetCacheDefinitions(t *testing.T) {
	p := &mainConfigurationFileProvider{}
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
		result := p.getCacheDefinitions(paths, caches)
		assert.Contains(t, result, "proxy_cache_path /var/cache/nginx/")
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
		result := p.getCacheDefinitions(paths, caches)
		assert.Contains(t, result, "proxy_cache_path /mnt/ssd/cache")
	})

	t.Run("generates basic config when optional fields are nil", func(t *testing.T) {
		caches := []cache.Cache{
			{
				ID: id1,
			},
		}
		result := p.getCacheDefinitions(paths, caches)
		assert.Contains(t, result, "proxy_cache_path /var/cache/nginx/")
		assert.NotContains(t, result, "inactive=")
		assert.NotContains(t, result, "max_size=")
	})
}
