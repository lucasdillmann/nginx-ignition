package cfgfiles

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/core/binding"
	"dillmann.com.br/nginx-ignition/core/cache"
	"dillmann.com.br/nginx-ignition/core/common/coreerror"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
	"dillmann.com.br/nginx-ignition/core/host"
	"dillmann.com.br/nginx-ignition/core/integration"
)

func Test_hostConfigurationFileProvider(t *testing.T) {
	t.Run("Provide", func(t *testing.T) {
		provider := &hostConfigurationFileProvider{}

		h := newHost()
		h.Routes = []host.Route{
			{
				Enabled:    true,
				Type:       host.ProxyRouteType,
				SourcePath: "/",
				TargetURI:  new("http://backend:8080"),
			},
		}

		ctx := newProviderContext(t)
		ctx.hosts = []host.Host{h}
		ctx.cfg = newSettings()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		integrationCmds := integration.NewMockedCommands(ctrl)
		provider.integrationCommands = integrationCmds

		files, err := provider.provide(ctx)
		assert.NoError(t, err)
		assert.Len(t, files, 1)
		assert.Equal(t, fmt.Sprintf("host-%s.conf", h.ID), files[0].Name)
		assert.Contains(t, files[0].Contents, "server_name _;")
		assert.Contains(t, files[0].Contents, "listen 0.0.0.0:80 default_server;")
		assert.Contains(t, files[0].Contents, "proxy_pass http://backend:8080;")
	})

	t.Run("Provide with traffic stats enabled", func(t *testing.T) {
		provider := &hostConfigurationFileProvider{}

		h := newHost()
		h.FeatureSet.StatsEnabled = true
		h.Routes = []host.Route{
			{
				Enabled:    true,
				Type:       host.ProxyRouteType,
				SourcePath: "/",
				TargetURI:  new("http://backend:8080"),
			},
		}

		ctx := newProviderContext(t)
		ctx.hosts = []host.Host{h}
		ctx.cfg = newSettings()
		ctx.cfg.Nginx.Stats.Enabled = true
		ctx.supportedFeatures.StatsType = StaticSupportType

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		integrationCmds := integration.NewMockedCommands(ctrl)
		provider.integrationCommands = integrationCmds

		files, err := provider.provide(ctx)
		assert.NoError(t, err)
		assert.Len(t, files, 1)
		assert.Contains(t, files[0].Contents, fmt.Sprintf("set $stats_host_id \"%s\";", h.ID))
		assert.Contains(t, files[0].Contents, "vhost_traffic_status on;")
		assert.Contains(
			t,
			files[0].Contents,
			"vhost_traffic_status_filter_by_set_key $stats_host_id hosts;",
		)
	})

	t.Run("Provide with global traffic stats enabled", func(t *testing.T) {
		provider := &hostConfigurationFileProvider{}

		h := newHost()
		h.FeatureSet.StatsEnabled = false
		h.Routes = []host.Route{
			{
				Enabled:    true,
				Type:       host.ProxyRouteType,
				SourcePath: "/",
				TargetURI:  new("http://backend:8080"),
			},
		}

		ctx := newProviderContext(t)
		ctx.hosts = []host.Host{h}
		ctx.cfg = newSettings()
		ctx.cfg.Nginx.Stats.Enabled = true
		ctx.cfg.Nginx.Stats.AllHosts = true
		ctx.supportedFeatures.StatsType = StaticSupportType

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		integrationCmds := integration.NewMockedCommands(ctrl)
		provider.integrationCommands = integrationCmds

		files, err := provider.provide(ctx)
		assert.NoError(t, err)
		assert.Len(t, files, 1)
		assert.Contains(t, files[0].Contents, "vhost_traffic_status on;")
	})

	t.Run("BuildServerNames", func(t *testing.T) {
		provider := &hostConfigurationFileProvider{}

		t.Run("returns underscore for default server", func(t *testing.T) {
			h := &host.Host{
				DefaultServer: true,
			}
			assert.Equal(t, "server_name _;", provider.buildServerNames(h))
		})

		t.Run("returns space separated domain names", func(t *testing.T) {
			h := &host.Host{
				DomainNames: []string{
					"example.com",
					"www.example.com",
				},
			}
			assert.Equal(
				t,
				"server_name example.com www.example.com;",
				provider.buildServerNames(h),
			)
		})
	})

	t.Run("BuildProxyPass", func(t *testing.T) {
		provider := &hostConfigurationFileProvider{}

		t.Run("returns simple proxy_pass", func(t *testing.T) {
			r := &host.Route{
				TargetURI: new("http://backend:8080"),
			}
			assert.Equal(t, "proxy_pass http://backend:8080;", provider.buildProxyPass(r))
		})

		t.Run("sets Host header when KeepOriginalDomainName is true", func(t *testing.T) {
			r := &host.Route{
				TargetURI: new("http://backend:8080"),
				Settings: host.RouteSettings{
					KeepOriginalDomainName: true,
				},
			}
			result := provider.buildProxyPass(r)
			assert.Contains(t, result, "proxy_pass http://backend:8080;")
			assert.Contains(t, result, "proxy_set_header Host backend:8080;")
		})

		t.Run("handles custom target URI override", func(t *testing.T) {
			r := &host.Route{
				TargetURI: new("http://default:8080"),
			}
			result := provider.buildProxyPass(r, "http://override:9090")
			assert.Equal(t, "proxy_pass http://override:9090;", result)
		})
	})

	t.Run("BuildRedirectRoute", func(t *testing.T) {
		provider := &hostConfigurationFileProvider{}
		ctx := newProviderContext(t)

		t.Run("generates redirect route config", func(t *testing.T) {
			r := &host.Route{
				SourcePath:   "/old",
				RedirectCode: new(301),
				TargetURI:    new("http://new.example.com"),
			}
			result := provider.buildRedirectRoute(ctx, r, host.FeatureSet{})
			assert.Contains(t, result, "location /old {")
			assert.Contains(t, result, "return 301 http://new.example.com;")
		})
	})

	t.Run("BuildIntegrationRoute", func(t *testing.T) {
		provider := &hostConfigurationFileProvider{}
		ctx := newProviderContext(t)

		t.Run("generates integration route config with dns resolvers", func(t *testing.T) {
			integrationID := uuid.New()
			r := &host.Route{
				SourcePath: "/api",
				Integration: &host.RouteIntegrationConfig{
					IntegrationID: integrationID,
					OptionID:      "opt-1",
				},
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			integrationCmds := integration.NewMockedCommands(ctrl)
			integrationCmds.EXPECT().
				GetOptionURL(gomock.Any(), integrationID, "opt-1").
				Return(new("http://1.2.3.4:80"), []string{"8.8.8.8", "8.8.4.4"}, nil)
			provider.integrationCommands = integrationCmds

			result, err := provider.buildIntegrationRoute(ctx, r, host.FeatureSet{})
			assert.NoError(t, err)
			assert.Contains(t, result, "location /api {")
			assert.Contains(t, result, "resolver 8.8.8.8 8.8.4.4 valid=5s;")
			assert.Contains(t, result, "proxy_pass http://1.2.3.4:80;")
		})

		t.Run("generates integration route config with target URI", func(t *testing.T) {
			integrationID := uuid.New()
			r := &host.Route{
				SourcePath: "/api",
				TargetURI:  new("/v1/resource"),
				Integration: &host.RouteIntegrationConfig{
					IntegrationID: integrationID,
					OptionID:      "opt-1",
				},
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			integrationCmds := integration.NewMockedCommands(ctrl)
			integrationCmds.EXPECT().
				GetOptionURL(gomock.Any(), integrationID, "opt-1").
				Return(new("http://1.2.3.4:80"), nil, nil)
			provider.integrationCommands = integrationCmds

			result, err := provider.buildIntegrationRoute(ctx, r, host.FeatureSet{})
			assert.NoError(t, err)
			assert.Contains(t, result, "proxy_pass http://1.2.3.4:80/v1/resource;")
		})

		t.Run("returns error when integration option not found", func(t *testing.T) {
			r := &host.Route{
				Integration: &host.RouteIntegrationConfig{},
			}
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			integrationCmds := integration.NewMockedCommands(ctrl)
			integrationCmds.EXPECT().
				GetOptionURL(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil, nil, nil)
			provider.integrationCommands = integrationCmds
			_, err := provider.buildIntegrationRoute(ctx, r, host.FeatureSet{})
			assert.Error(t, err)
			var coreErr *coreerror.CoreError
			assert.ErrorAs(t, err, &coreErr)
			assert.Equal(t, i18n.K.CoreNginxCfgfilesOptionNotFound, coreErr.Message.Key)
		})
	})

	t.Run("BuildExecuteCodeRoute", func(t *testing.T) {
		provider := &hostConfigurationFileProvider{}
		ctx := newProviderContext(t)
		h := &host.Host{ID: uuid.New()}

		t.Run("generates javascript route config", func(t *testing.T) {
			r := &host.Route{
				Priority:   1,
				SourcePath: "/js",
				SourceCode: &host.RouteSourceCode{
					Language:     host.JavascriptCodeLanguage,
					MainFunction: new("handler"),
				},
			}
			result, err := provider.buildExecuteCodeRoute(ctx, h, r)
			assert.NoError(t, err)
			assert.Contains(
				t,
				result,
				fmt.Sprintf("js_import route_1 from \"/etc/nginx/host-%s-route-1.js\";", h.ID),
			)
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
			result, err := provider.buildExecuteCodeRoute(ctx, h, r)
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
			_, err := provider.buildExecuteCodeRoute(ctx, h, r)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "invalid language")
		})
	})

	t.Run("BuildStaticResponseRoute", func(t *testing.T) {
		provider := &hostConfigurationFileProvider{}
		ctx := newProviderContext(t)
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
			result := provider.buildStaticResponseRoute(ctx, h, r)
			assert.Contains(t, result, "location @route_2/static_payload {")
			assert.Contains(t, result, "add_header \"Content-Type\" \"application/json\" always;")
			assert.Contains(
				t,
				result,
				"try_files \"/host-"+h.ID.String()+"-route-2.payload\" =200;",
			)
		})
	})

	t.Run("BuildRouteFeatures", func(t *testing.T) {
		provider := &hostConfigurationFileProvider{}

		t.Run("returns websocket config when enabled", func(t *testing.T) {
			features := host.FeatureSet{
				WebsocketSupport: true,
			}
			result := provider.buildRouteFeatures(features)
			assert.Contains(t, result, "proxy_http_version 1.1;")
			assert.Contains(t, result, "proxy_set_header Upgrade $http_upgrade;")
			assert.Contains(t, result, "proxy_set_header Connection \"upgrade\";")
		})

		t.Run("returns empty string when disabled", func(t *testing.T) {
			features := host.FeatureSet{
				WebsocketSupport: false,
			}
			assert.Equal(t, "", provider.buildRouteFeatures(features))
		})
	})

	t.Run("BuildRouteSettings", func(t *testing.T) {
		provider := &hostConfigurationFileProvider{}
		ctx := newProviderContext(t)

		t.Run("includes forward headers when enabled", func(t *testing.T) {
			r := &host.Route{
				Settings: host.RouteSettings{
					IncludeForwardHeaders: true,
				},
			}
			result := provider.buildRouteSettings(ctx, r)
			assert.Contains(
				t,
				result,
				"proxy_set_header x-forwarded-for $proxy_add_x_forwarded_for;",
			)
			assert.Contains(t, result, "proxy_set_header x-real-ip $remote_addr;")
		})

		t.Run("includes custom configuration", func(t *testing.T) {
			r := &host.Route{
				Settings: host.RouteSettings{
					Custom: new("proxy_buffer_size 16k;"),
				},
			}
			result := provider.buildRouteSettings(ctx, r)
			assert.Contains(t, result, "proxy_buffer_size 16k;")
		})

		t.Run("includes access list when present", func(t *testing.T) {
			id := uuid.New()
			r := &host.Route{
				AccessListID: &id,
			}
			result := provider.buildRouteSettings(ctx, r)
			assert.Contains(
				t,
				result,
				fmt.Sprintf("include \"/etc/nginx/access-list-%s.conf\";", id),
			)
		})
	})

	t.Run("BuildBinding", func(t *testing.T) {
		provider := &hostConfigurationFileProvider{}
		ctx := newProviderContext(t)
		ctx.cfg = newSettings()
		h := &host.Host{ID: uuid.New()}
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		t.Run("generates HTTP binding", func(t *testing.T) {
			b := &binding.Binding{
				Type: binding.HTTPBindingType,
				IP:   "127.0.0.1",
				Port: 8080,
			}
			result, err := provider.buildBinding(
				ctx,
				h,
				b,
				[]string{},
				"server_name example.com;",
				"",
				"",
				"",
			)
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
			result, err := provider.buildBinding(
				ctx,
				h,
				b,
				[]string{},
				"server_name example.com;",
				"",
				"",
				"",
			)
			assert.NoError(t, err)
			assert.Contains(t, result, "listen 0.0.0.0:443 ssl ;")
			assert.Contains(
				t,
				result,
				fmt.Sprintf("ssl_certificate \"/etc/nginx/certificate-%s.pem\";", certID),
			)
		})

		t.Run("includes HTTP to HTTPS redirect in HTTP binding", func(t *testing.T) {
			b := &binding.Binding{Type: binding.HTTPBindingType}
			redirect := "return 301 https://$server_name$request_uri;"
			result, err := provider.buildBinding(ctx, h, b, []string{}, "", redirect, "", "")
			assert.NoError(t, err)
			assert.Contains(t, result, redirect)
		})

		t.Run("includes HTTP2 in HTTPS binding", func(t *testing.T) {
			b := &binding.Binding{
				Type:          binding.HTTPSBindingType,
				CertificateID: new(uuid.New()),
			}
			result, err := provider.buildBinding(ctx, h, b, []string{}, "", "", "http2 on;", "")
			assert.NoError(t, err)
			assert.Contains(t, result, "http2 on;")
		})

		t.Run("returns error for invalid binding type", func(t *testing.T) {
			b := &binding.Binding{Type: "INVALID"}
			_, err := provider.buildBinding(ctx, h, b, []string{}, "", "", "", "")
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "invalid binding type")
		})
	})

	t.Run("BuildRoute", func(t *testing.T) {
		provider := &hostConfigurationFileProvider{}
		ctx := newProviderContext(t)
		h := &host.Host{}

		t.Run("returns error for invalid route type", func(t *testing.T) {
			r := &host.Route{Type: "INVALID"}
			_, err := provider.buildRoute(ctx, h, r)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "invalid route type")
		})
	})

	t.Run("BuildCacheConfig", func(t *testing.T) {
		provider := &hostConfigurationFileProvider{}
		cacheID := uuid.New()
		c := newCache()
		c.ID = cacheID
		c.MinimumUsesBeforeCaching = 2
		c.BackgroundUpdate = true
		c.Revalidate = true
		c.Durations = []cache.Duration{
			{
				StatusCodes:      []string{"200", "302"},
				ValidTimeSeconds: 600,
			},
		}
		c.AllowedMethods = []cache.Method{
			cache.GetMethod,
			cache.HeadMethod,
		}
		c.IgnoreUpstreamCacheHeaders = true
		c.CacheStatusResponseHeaderEnabled = true
		c.UseStale = []cache.UseStaleOption{
			cache.ErrorUseStale,
			cache.TimeoutUseStale,
		}
		c.ConcurrencyLock = cache.ConcurrencyLock{
			Enabled:        true,
			TimeoutSeconds: new(5),
			AgeSeconds:     new(10),
		}
		c.BypassRules = []string{"$cookie_nocache"}
		c.NoCacheRules = []string{"$arg_nocache"}
		c.FileExtensions = []string{"jpg", "png"}

		caches := []cache.Cache{c}

		t.Run("generates comprehensive cache config", func(t *testing.T) {
			result := provider.buildCacheConfig(caches, &cacheID)
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
			result := provider.buildCacheConfig(caches, new(uuid.New()))
			assert.Equal(t, "", result)
		})

		t.Run("returns empty string when cacheID is nil", func(t *testing.T) {
			result := provider.buildCacheConfig(caches, nil)
			assert.Equal(t, "", result)
		})
	})

	t.Run("BuildStaticFilesRoute", func(t *testing.T) {
		provider := &hostConfigurationFileProvider{}
		ctx := newProviderContext(t)

		t.Run("generates static files config", func(t *testing.T) {
			r := &host.Route{
				SourcePath: "/static",
				TargetURI:  new("/var/www/static"),
				Settings: host.RouteSettings{
					DirectoryListingEnabled: true,
				},
			}
			result := provider.buildStaticFilesRoute(ctx, r)
			assert.Contains(t, result, "location /static/ {")
			assert.Contains(t, result, "root \"/var/www/static\";")
			assert.Contains(t, result, "autoindex on;")
		})

		t.Run("generates static files config with index file", func(t *testing.T) {
			r := &host.Route{
				SourcePath: "/static",
				TargetURI:  new("/var/www/static"),
				Settings: host.RouteSettings{
					IndexFile: new("home.html"),
				},
			}
			result := provider.buildStaticFilesRoute(ctx, r)
			assert.Contains(t, result, "index \"home.html\";")
		})
	})
}
