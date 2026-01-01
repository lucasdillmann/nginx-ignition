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

func TestHostConfigurationFileProvider_BuildRedirectRoute(t *testing.T) {
	p := &hostConfigurationFileProvider{}
	ctx := &providerContext{}

	t.Run("generates redirect route config", func(t *testing.T) {
		r := &host.Route{
			SourcePath:   "/old",
			RedirectCode: ptr.Of(301),
			TargetURI:    ptr.Of("http://new.example.com"),
		}
		result := p.buildRedirectRoute(ctx, r, host.FeatureSet{})
		assert.Contains(t, result, "location /old {")
		assert.Contains(t, result, "return 301 http://new.example.com;")
	})
}

func TestHostConfigurationFileProvider_BuildIntegrationRoute(t *testing.T) {
	p := &hostConfigurationFileProvider{}
	ctx := &providerContext{
		context: context.Background(),
	}

	t.Run("generates integration route config with dns resolvers", func(t *testing.T) {
		integrationID := uuid.New()
		r := &host.Route{
			SourcePath: "/api",
			Integration: &host.RouteIntegrationConfig{
				IntegrationID: integrationID,
				OptionID:      "opt-1",
			},
		}

		p.integrationCommands = &integration.Commands{
			GetOptionURL: func(_ context.Context, iID uuid.UUID, oID string) (*string, []string, error) {
				assert.Equal(t, integrationID, iID)
				assert.Equal(t, "opt-1", oID)
				return ptr.Of("http://1.2.3.4:80"), []string{"8.8.8.8", "8.8.4.4"}, nil
			},
		}

		result, err := p.buildIntegrationRoute(ctx, r, host.FeatureSet{})
		assert.NoError(t, err)
		assert.Contains(t, result, "location /api {")
		assert.Contains(t, result, "resolver 8.8.8.8 8.8.4.4 valid=5s;")
		assert.Contains(t, result, "proxy_pass http://1.2.3.4:80;")
	})

	t.Run("returns error when integration option not found", func(t *testing.T) {
		r := &host.Route{
			Integration: &host.RouteIntegrationConfig{},
		}
		p.integrationCommands = &integration.Commands{
			GetOptionURL: func(_ context.Context, _ uuid.UUID, _ string) (*string, []string, error) {
				return nil, nil, nil
			},
		}
		_, err := p.buildIntegrationRoute(ctx, r, host.FeatureSet{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Integration option not found")
	})
}

func TestHostConfigurationFileProvider_BuildExecuteCodeRoute(t *testing.T) {
	p := &hostConfigurationFileProvider{}
	ctx := &providerContext{
		paths: &Paths{Config: "/etc/nginx/"},
	}
	h := &host.Host{ID: uuid.New()}

	t.Run("generates javascript route config", func(t *testing.T) {
		r := &host.Route{
			Priority:   1,
			SourcePath: "/js",
			SourceCode: &host.RouteSourceCode{
				Language:     host.JavascriptCodeLanguage,
				MainFunction: ptr.Of("handler"),
			},
		}
		result, err := p.buildExecuteCodeRoute(ctx, h, r)
		assert.NoError(t, err)
		assert.Contains(t, result, fmt.Sprintf("js_import route_1 from /etc/nginx/host-%s-route-1.js;", h.ID))
		assert.Contains(t, result, "js_content route_1.handler;")
	})

	t.Run("generates lua route config", func(t *testing.T) {
		r := &host.Route{
			SourcePath: "/lua",
			SourceCode: &host.RouteSourceCode{
				Language: host.LuaCodeLanguage,
				Contents: "ngx.say('hello')",
			},
		}
		result, err := p.buildExecuteCodeRoute(ctx, h, r)
		assert.NoError(t, err)
		assert.Contains(t, result, "content_by_lua_block")
		assert.Contains(t, result, "ngx.say('hello')")
	})

	t.Run("returns error for invalid language", func(t *testing.T) {
		r := &host.Route{
			SourceCode: &host.RouteSourceCode{
				Language: "FORTRAN",
			},
		}
		_, err := p.buildExecuteCodeRoute(ctx, h, r)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid language")
	})
}

func TestHostConfigurationFileProvider_BuildStaticResponseRoute(t *testing.T) {
	p := &hostConfigurationFileProvider{}
	ctx := &providerContext{
		paths: &Paths{Config: "/etc/nginx/"},
	}
	h := &host.Host{ID: uuid.New()}

	t.Run("generates static response route config", func(t *testing.T) {
		r := &host.Route{
			Priority:   2,
			SourcePath: "/static",
			Response: &host.RouteStaticResponse{
				StatusCode: 200,
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
			},
		}
		result := p.buildStaticResponseRoute(ctx, h, r)
		assert.Contains(t, result, "location @route_2/static_payload {")
		assert.Contains(t, result, "add_header \"Content-Type\" \"application/json\" always;")
		assert.Contains(t, result, "try_files /host-"+h.ID.String()+"-route-2.payload =200;")
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

func TestHostConfigurationFileProvider_BuildBinding(t *testing.T) {
	p := &hostConfigurationFileProvider{}
	paths := &Paths{
		Config: "/etc/nginx/",
		Logs:   "/var/log/nginx/",
	}
	ctx := &providerContext{
		context: context.Background(),
		paths:   paths,
	}
	h := &host.Host{ID: uuid.New()}
	p.settingsCommands = &settings.Commands{
		Get: func(_ context.Context) (*settings.Settings, error) {
			return &settings.Settings{
				Nginx: &settings.NginxSettings{
					Logs: &settings.NginxLogsSettings{
						AccessLogsEnabled: true,
						ErrorLogsEnabled:  true,
						ErrorLogsLevel:    settings.WarnLogLevel,
					},
				},
			}, nil
		},
	}

	t.Run("generates HTTP binding", func(t *testing.T) {
		b := &binding.Binding{
			Type: binding.HTTPBindingType,
			IP:   "127.0.0.1",
			Port: 8080,
		}
		result, err := p.buildBinding(ctx, h, b, []string{}, "server_name example.com;", "", "")
		assert.NoError(t, err)
		assert.Contains(t, result, "listen 127.0.0.1:8080 ;")
	})

	t.Run("generates HTTPS binding", func(t *testing.T) {
		certID := uuid.New()
		b := &binding.Binding{
			Type:          binding.HTTPSBindingType,
			IP:            "0.0.0.0",
			Port:          443,
			CertificateID: &certID,
		}
		result, err := p.buildBinding(ctx, h, b, []string{}, "server_name example.com;", "", "")
		assert.NoError(t, err)
		assert.Contains(t, result, "listen 0.0.0.0:443 ssl ;")
		assert.Contains(t, result, fmt.Sprintf("ssl_certificate /etc/nginx/certificate-%s.pem;", certID))
	})

	t.Run("includes HTTP to HTTPS redirect in HTTP binding", func(t *testing.T) {
		b := &binding.Binding{Type: binding.HTTPBindingType}
		redirect := "return 301 https://$server_name$request_uri;"
		result, err := p.buildBinding(ctx, h, b, []string{}, "", redirect, "")
		assert.NoError(t, err)
		assert.Contains(t, result, redirect)
	})

	t.Run("includes HTTP2 in HTTPS binding", func(t *testing.T) {
		certID := uuid.New()
		b := &binding.Binding{
			Type:          binding.HTTPSBindingType,
			CertificateID: &certID,
		}
		result, err := p.buildBinding(ctx, h, b, []string{}, "", "", "http2 on;")
		assert.NoError(t, err)
		assert.Contains(t, result, "http2 on;")
	})

	t.Run("returns error for invalid binding type", func(t *testing.T) {
		b := &binding.Binding{Type: "INVALID"}
		_, err := p.buildBinding(ctx, h, b, []string{}, "", "", "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid binding type")
	})

	t.Run("returns error when settingsCommands fails", func(t *testing.T) {
		p.settingsCommands = &settings.Commands{
			Get: func(_ context.Context) (*settings.Settings, error) {
				return nil, assert.AnError
			},
		}
		b := &binding.Binding{Type: binding.HTTPBindingType}
		_, err := p.buildBinding(ctx, h, b, []string{}, "", "", "")
		assert.ErrorIs(t, err, assert.AnError)
	})
}

func TestHostConfigurationFileProvider_BuildRoute(t *testing.T) {
	p := &hostConfigurationFileProvider{}
	ctx := &providerContext{}
	h := &host.Host{}

	t.Run("returns error for invalid route type", func(t *testing.T) {
		r := &host.Route{Type: "INVALID"}
		_, err := p.buildRoute(ctx, h, r)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid route type")
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
			IgnoreUpstreamCacheHeaders:       true,
			CacheStatusResponseHeaderEnabled: true,
			UseStale: []cache.UseStaleOption{
				cache.ErrorUseStale,
				cache.TimeoutUseStale,
			},
			ConcurrencyLock: cache.ConcurrencyLock{
				Enabled:        true,
				TimeoutSeconds: ptr.Of(5),
				AgeSeconds:     ptr.Of(10),
			},
			BypassRules:    []string{"$cookie_nocache"},
			NoCacheRules:   []string{"$arg_nocache"},
			FileExtensions: []string{"jpg", "png"},
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
		assert.Contains(t, result, "proxy_ignore_headers Cache-Control Expires;")
		assert.Contains(t, result, "add_header X-Cache-Status $upstream_cache_status;")
		assert.Contains(t, result, "proxy_cache_use_stale error timeout;")
		assert.Contains(t, result, "proxy_cache_lock on;")
		assert.Contains(t, result, "proxy_cache_lock_timeout 5s;")
		assert.Contains(t, result, "proxy_cache_lock_age 10s;")
		assert.Contains(t, result, "proxy_cache_bypass $cookie_nocache;")
		assert.Contains(t, result, "proxy_no_cache $arg_nocache;")
		assert.Contains(t, result, "if ($uri !~* \"\\.(jpg|png)$\")")
	})

	t.Run("returns empty string when cache not found", func(t *testing.T) {
		unknownID := uuid.New()
		result := p.buildCacheConfig(caches, &unknownID)
		assert.Equal(t, "", result)
	})

	t.Run("returns empty string when cacheID is nil", func(t *testing.T) {
		result := p.buildCacheConfig(caches, nil)
		assert.Equal(t, "", result)
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
