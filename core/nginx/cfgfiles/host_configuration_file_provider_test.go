package cfgfiles

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"dillmann.com.br/nginx-ignition/core/binding"
	"dillmann.com.br/nginx-ignition/core/cache"
	"dillmann.com.br/nginx-ignition/core/common/ptr"
	"dillmann.com.br/nginx-ignition/core/host"
	"dillmann.com.br/nginx-ignition/core/integration"
	"dillmann.com.br/nginx-ignition/core/settings"
)

func TestHostConfigurationFileProvider_Provide(t *testing.T) {
	p := &hostConfigurationFileProvider{}
	paths := &Paths{
		Config: "/etc/nginx/",
		Logs:   "/var/log/nginx/",
	}
	id := uuid.New()
	ctx := &providerContext{
		context: context.Background(),
		paths:   paths,
		hosts: []host.Host{
			{
				ID:            id,
				Enabled:       true,
				DefaultServer: true,
				DomainNames:   []string{"example.com"},
				Bindings: []binding.Binding{
					{
						Type: binding.HTTPBindingType,
						IP:   "0.0.0.0",
						Port: 80,
					},
				},
				Routes: []host.Route{
					{
						Enabled:    true,
						Type:       host.ProxyRouteType,
						SourcePath: "/",
						TargetURI:  ptr.Of("http://backend:8080"),
					},
				},
			},
		},
	}

	p.settingsCommands = &settings.Commands{
		Get: func(_ context.Context) (*settings.Settings, error) {
			return &settings.Settings{
				Nginx: &settings.NginxSettings{
					WorkerProcesses: 1,
					Logs: &settings.NginxLogsSettings{
						AccessLogsEnabled: true,
						ErrorLogsEnabled:  true,
						ErrorLogsLevel:    settings.WarnLogLevel,
					},
				},
			}, nil
		},
	}

	p.integrationCommands = &integration.Commands{}

	files, err := p.provide(ctx)
	assert.NoError(t, err)
	assert.Len(t, files, 1)
	assert.Equal(t, fmt.Sprintf("host-%s.conf", id), files[0].Name)
	assert.Contains(t, files[0].Contents, "server_name _;")
	assert.Contains(t, files[0].Contents, "listen 0.0.0.0:80 default_server;")
	assert.Contains(t, files[0].Contents, "proxy_pass http://backend:8080;")
}

func TestHostConfigurationFileProvider_BuildServerNames(t *testing.T) {
	p := &hostConfigurationFileProvider{}

	t.Run("returns underscore for default server", func(t *testing.T) {
		h := &host.Host{
			DefaultServer: true,
		}
		assert.Equal(t, "server_name _;", p.buildServerNames(h))
	})

	t.Run("returns space separated domain names", func(t *testing.T) {
		h := &host.Host{
			DomainNames: []string{
				"example.com",
				"www.example.com",
			},
		}
		assert.Equal(t, "server_name example.com www.example.com;", p.buildServerNames(h))
	})
}

func TestHostConfigurationFileProvider_BuildProxyPass(t *testing.T) {
	p := &hostConfigurationFileProvider{}

	t.Run("returns simple proxy_pass", func(t *testing.T) {
		r := &host.Route{
			TargetURI: ptr.Of("http://backend:8080"),
		}
		assert.Equal(t, "proxy_pass http://backend:8080;", p.buildProxyPass(r))
	})

	t.Run("sets Host header when KeepOriginalDomainName is true", func(t *testing.T) {
		r := &host.Route{
			TargetURI: ptr.Of("http://backend:8080"),
			Settings: host.RouteSettings{
				KeepOriginalDomainName: true,
			},
		}
		result := p.buildProxyPass(r)
		assert.Contains(t, result, "proxy_pass http://backend:8080;")
		assert.Contains(t, result, "proxy_set_header Host backend:8080;")
	})

	t.Run("handles custom target URI override", func(t *testing.T) {
		r := &host.Route{
			TargetURI: ptr.Of("http://default:8080"),
		}
		result := p.buildProxyPass(r, "http://override:9090")
		assert.Equal(t, "proxy_pass http://override:9090;", result)
	})
}

func TestHostConfigurationFileProvider_BuildRouteFeatures(t *testing.T) {
	p := &hostConfigurationFileProvider{}

	t.Run("returns websocket config when enabled", func(t *testing.T) {
		features := host.FeatureSet{
			WebsocketSupport: true,
		}
		result := p.buildRouteFeatures(features)
		assert.Contains(t, result, "proxy_http_version 1.1;")
		assert.Contains(t, result, "proxy_set_header Upgrade $http_upgrade;")
		assert.Contains(t, result, "proxy_set_header Connection \"upgrade\";")
	})

	t.Run("returns empty string when disabled", func(t *testing.T) {
		features := host.FeatureSet{
			WebsocketSupport: false,
		}
		assert.Equal(t, "", p.buildRouteFeatures(features))
	})
}

func TestHostConfigurationFileProvider_BuildRouteSettings(t *testing.T) {
	p := &hostConfigurationFileProvider{}
	ctx := &providerContext{
		paths: &Paths{
			Config: "/etc/nginx/",
		},
	}

	t.Run("includes forward headers when enabled", func(t *testing.T) {
		r := &host.Route{
			Settings: host.RouteSettings{
				IncludeForwardHeaders: true,
			},
		}
		result := p.buildRouteSettings(ctx, r)
		assert.Contains(t, result, "proxy_set_header x-forwarded-for $proxy_add_x_forwarded_for;")
		assert.Contains(t, result, "proxy_set_header x-real-ip $remote_addr;")
	})

	t.Run("includes custom configuration", func(t *testing.T) {
		r := &host.Route{
			Settings: host.RouteSettings{
				Custom: ptr.Of("proxy_buffer_size 16k;"),
			},
		}
		result := p.buildRouteSettings(ctx, r)
		assert.Contains(t, result, "proxy_buffer_size 16k;")
	})

	t.Run("includes access list when present", func(t *testing.T) {
		id := uuid.New()
		r := &host.Route{
			AccessListID: &id,
		}
		result := p.buildRouteSettings(ctx, r)
		assert.Contains(t, result, fmt.Sprintf("include /etc/nginx/access-list-%s.conf;", id))
	})
}

func TestHostConfigurationFileProvider_BuildCacheConfig(t *testing.T) {
	p := &hostConfigurationFileProvider{}
	cacheID := uuid.New()
	caches := []cache.Cache{
		{
			ID:                       cacheID,
			MinimumUsesBeforeCaching: 2,
			BackgroundUpdate:         true,
			Revalidate:               true,
			Durations: []cache.Duration{
				{
					StatusCodes:      []string{"200", "302"},
					ValidTimeSeconds: 600,
				},
			},
			AllowedMethods: []cache.Method{
				cache.GetMethod,
				cache.HeadMethod,
			},
		},
	}

	t.Run("generates comprehensive cache config", func(t *testing.T) {
		result := p.buildCacheConfig(caches, &cacheID)
		cacheIDNoDashes := strings.ReplaceAll(cacheID.String(), "-", "")
		assert.Contains(t, result, fmt.Sprintf("proxy_cache cache_%s;", cacheIDNoDashes))
		assert.Contains(t, result, "proxy_cache_min_uses 2;")
		assert.Contains(t, result, "proxy_cache_background_update on;")
		assert.Contains(t, result, "proxy_cache_revalidate on;")
		assert.Contains(t, result, "proxy_cache_valid 200 302 600s;")
		assert.Contains(t, result, "proxy_cache_methods get head;")
	})
}

func TestHostConfigurationFileProvider_BuildStaticFilesRoute(t *testing.T) {
	p := &hostConfigurationFileProvider{}
	ctx := &providerContext{}

	t.Run("generates static files config", func(t *testing.T) {
		r := &host.Route{
			SourcePath: "/static",
			TargetURI:  ptr.Of("/var/www/static"),
			Settings: host.RouteSettings{
				DirectoryListingEnabled: true,
			},
		}
		result := p.buildStaticFilesRoute(ctx, r)
		assert.Contains(t, result, "location /static/ {")
		assert.Contains(t, result, "root /var/www/static;")
		assert.Contains(t, result, "autoindex on;")
	})
}
