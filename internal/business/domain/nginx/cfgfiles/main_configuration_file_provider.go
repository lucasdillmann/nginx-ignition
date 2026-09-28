package cfgfiles

import (
	"path/filepath"
	"strings"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/configuration"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/runtime"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/cache"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/host"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/settings"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/stream"
)

type mainConfigurationFileProvider struct {
	config *configuration.Configuration
}

func newMainConfigurationFileProvider(
	config *configuration.Configuration,
) *mainConfigurationFileProvider {
	return &mainConfigurationFileProvider{
		config: config,
	}
}

func (p *mainConfigurationFileProvider) provide(ctx *providerContext) ([]File, error) {
	cfg := ctx.cfg
	logs := cfg.Nginx.Logs
	moduleLines := strings.Builder{}
	streamLines := strings.Builder{}

	if ctx.supportedFeatures.StreamType != NoneSupportType {
		if ctx.supportedFeatures.StreamType == DynamicSupportType {
			_, _ = moduleLines.WriteString("load_module modules/ngx_stream_module.so;\n")
		}

		_, _ = streamLines.WriteString("stream {\n")
		_, _ = streamLines.WriteString(p.getStreamIncludes(ctx.paths, ctx.streams))
		_, _ = streamLines.WriteString("}\n")
	}

	if ctx.supportedFeatures.RunCodeType == DynamicSupportType {
		_, _ = moduleLines.WriteString("load_module modules/ndk_http_module.so;\n")
		_, _ = moduleLines.WriteString("load_module modules/ngx_http_js_module.so;\n")
		_, _ = moduleLines.WriteString("load_module modules/ngx_http_lua_module.so;\n")
	}

	if ctx.supportedFeatures.StatsType == DynamicSupportType {
		_, _ = moduleLines.WriteString("load_module modules/ngx_http_geoip2_module.so;\n")
		_, _ = moduleLines.WriteString(
			"load_module modules/ngx_http_vhost_traffic_status_module.so;\n",
		)
	}

	var customCfg string
	if cfg.Nginx.Custom != nil {
		customCfg = nginxSprintf("\n%s\n", rawConfigFragment(*cfg.Nginx.Custom))
	}

	userStatement := nginxSprintf("user %s %s;", cfg.Nginx.RuntimeUser, cfg.Nginx.RuntimeUser)
	if runtime.IsWindows() {
		userStatement = ""
	}

	statsDefinitions, err := p.getStatsDefinitions(ctx.paths, cfg.Nginx.Stats)
	if err != nil {
		return nil, err
	}

	contents := nginxSprintf(
		`
			%s
			%s
			worker_processes %d;
			pid %s;
			error_log %s;
			
			events {
				worker_connections %d;
			}
			
			http {
				sendfile %s;
				server_tokens %s;
				tcp_nodelay %s;
				
				keepalive_timeout %ds;
				proxy_connect_timeout %ds;
				proxy_read_timeout %ds;
				proxy_send_timeout %ds;
				send_timeout %ds;
				client_body_timeout %ds;

				client_max_body_size %dM;
				client_body_buffer_size %dk;
				client_header_buffer_size %dk;
				large_client_header_buffers %d %dk;
				output_buffers %d %dk;
				client_body_temp_path %s 1 2;
				proxy_temp_path %s 1 2;
				fastcgi_temp_path %s 1 2;
				scgi_temp_path %s 1 2;
				uwsgi_temp_path %s 1 2;

				default_type %s;
				include %s;
				%s
				%s
				%s
				%s
			}
			
			%s
		`,
		directiveFragment(userStatement),
		directiveFragment(moduleLines.String()),
		cfg.Nginx.WorkerProcesses,
		ctx.paths.Base+"nginx.pid",
		directiveFragment(p.getErrorLogPath(ctx.paths, logs)),
		cfg.Nginx.WorkerConnections,
		directiveFragment(statusFlag(cfg.Nginx.SendfileEnabled)),
		directiveFragment(statusFlag(cfg.Nginx.ServerTokensEnabled)),
		directiveFragment(statusFlag(cfg.Nginx.TCPNoDelayEnabled)),
		cfg.Nginx.Timeouts.Keepalive,
		cfg.Nginx.Timeouts.Connect,
		cfg.Nginx.Timeouts.Read,
		cfg.Nginx.Timeouts.Send,
		cfg.Nginx.Timeouts.Send,
		cfg.Nginx.Timeouts.ClientBody,
		cfg.Nginx.MaximumBodySizeMb,
		cfg.Nginx.Buffers.ClientBodyKb,
		cfg.Nginx.Buffers.ClientHeaderKb,
		cfg.Nginx.Buffers.LargeClientHeader.Amount,
		cfg.Nginx.Buffers.LargeClientHeader.SizeKb,
		cfg.Nginx.Buffers.Output.Amount,
		cfg.Nginx.Buffers.Output.SizeKb,
		filepath.ToSlash(filepath.Join(ctx.paths.Temp, "client-body")),
		filepath.ToSlash(filepath.Join(ctx.paths.Temp, "proxy")),
		filepath.ToSlash(filepath.Join(ctx.paths.Temp, "fastcgi")),
		filepath.ToSlash(filepath.Join(ctx.paths.Temp, "scgi")),
		filepath.ToSlash(filepath.Join(ctx.paths.Temp, "uwsgi")),
		cfg.Nginx.DefaultContentType,
		ctx.paths.Config+"mime.types",
		rawConfigFragment(customCfg),
		directiveFragment(p.getCacheDefinitions(ctx.paths, ctx.caches)),
		directiveFragment(statsDefinitions),
		directiveFragment(p.getHostIncludes(ctx.paths, ctx.hosts)),
		directiveFragment(streamLines.String()),
	)

	return []File{
		{
			Name:     "nginx.conf",
			Contents: contents,
		},
	}, nil
}

func (p *mainConfigurationFileProvider) getErrorLogPath(
	paths *Paths,
	logs *settings.NginxLogsSettings,
) string {
	if logs.ServerLogsEnabled {
		return nginxSprintf(
			"%s %s",
			paths.Logs+"main.log",
			strings.ToLower(string(logs.ServerLogsLevel)),
		)
	}

	return "off"
}

func (p *mainConfigurationFileProvider) getHostIncludes(paths *Paths, hosts []host.Host) string {
	includes := make([]string, 0, len(hosts))
	for _, h := range hosts {
		includes = append(
			includes,
			nginxSprintf("include %s;", paths.Config+"host-"+h.ID.String()+".conf"),
		)
	}

	return strings.Join(includes, "\n")
}

func (p *mainConfigurationFileProvider) getStreamIncludes(
	paths *Paths,
	streams []stream.Stream,
) string {
	includes := make([]string, 0, len(streams))

	for _, s := range streams {
		includes = append(
			includes,
			nginxSprintf("include %s;", paths.Config+"stream-"+s.ID.String()+".conf"),
		)
	}

	return strings.Join(includes, "\n")
}

func (p *mainConfigurationFileProvider) getCacheDefinitions(
	paths *Paths,
	caches []cache.Cache,
) string {
	if len(caches) == 0 {
		return ""
	}

	results := make([]string, 0)
	for _, c := range caches {
		cacheIDNoDashes := strings.ReplaceAll(c.ID.String(), "-", "")
		storagePath := c.StoragePath

		if storagePath == nil || strings.TrimSpace(*storagePath) == "" {
			storagePath = new(paths.Cache + cacheIDNoDashes)
		}

		inactive := ""
		if c.InactiveSeconds != nil {
			inactive = " " + nginxSprintf("inactive=%ds", *c.InactiveSeconds)
		}

		maxSize := ""
		if c.MaximumSizeMB != nil {
			maxSize = " " + nginxSprintf("max_size=%dm", *c.MaximumSizeMB)
		}

		results = append(results, nginxSprintf(
			"proxy_cache_path %s levels=1:2 keys_zone=%s:10m%s%s;",
			*storagePath,
			directiveFragment("cache_"+cacheIDNoDashes),
			directiveFragment(inactive),
			directiveFragment(maxSize),
		))
	}

	return strings.Join(results, "\n")
}

func (p *mainConfigurationFileProvider) getStatsDefinitions(
	paths *Paths,
	cfg *settings.NginxStatsSettings,
) (string, error) {
	if cfg == nil || !cfg.Enabled {
		return "", nil
	}

	geoIPCountryFilePath := filepath.Join(paths.Config, "geoip-country.mmdb")
	geoIPCityFilePath := filepath.Join(paths.Config, "geoip-city.mmdb")
	output := strings.Builder{}

	nginxFprintf(
		&output,
		`
		geoip2 %s {
			$geoip_country_code default=Unknown source=$remote_addr country iso_code;
		}

		geoip2 %s {
			$geoip_city_name default=Unknown source=$remote_addr city names en;
		}
		
		map $http_user_agent $stats_user_agent {
			default "Unknown";
			
			"~(?i)bot|crawler|spider|slurp|google|bing|yandex|baidu|duckduckbot|twitterbot|facebookexternalhit|linkedinbot|slackbot|wget|curl" "Bots";
			"~(?i)playstation|nintendo|xbox" "Console";
			"~(?i)smart-tv|smarttv|googletv|appletv|hbbtv|tizen|samsungview|crkey" "Smart TV";
			"~(?i)macintosh|mac\sos\sx" "macOS";
			"~(?i)ipad" "iPadOS";
			"~(?i)iphone|ipod" "iOS";
			"~(?i)android" "Android";
			"~(?i)windows" "Windows";
			"~(?i)linux" "Linux";
		}
		
		vhost_traffic_status_zone shared:nginx-ignition-traffic-stats:%dm;
		vhost_traffic_status_filter_by_host on;
		vhost_traffic_status_stats_by_upstream on;
		vhost_traffic_status_filter_by_set_key $geoip_country_code countryCode@global;
		vhost_traffic_status_filter_by_set_key $geoip_city_name city@global;
		vhost_traffic_status_filter_by_set_key $stats_user_agent userAgent@global;
		`,
		geoIPCountryFilePath,
		geoIPCityFilePath,
		cfg.MaximumSizeMB,
	)

	if cfg.Persistent {
		dbLocation := cfg.DatabaseLocation
		if dbLocation == nil {
			dataPath, err := p.config.Get("nginx-ignition.database.data-path")
			if err != nil {
				return "", err
			}

			dbLocation = new(filepath.Join(dataPath, "traffic-stats.db"))
		}

		nginxFprintf(&output, "vhost_traffic_status_dump %s 5s;\n", *dbLocation)
	}

	nginxFprintf(
		&output,
		`
		server {
			root /dev/null;
            listen %s;
			access_log off;
			vhost_traffic_status off;
			
			location / {
				vhost_traffic_status_display;
				vhost_traffic_status_display_format json;
			}
        }
		`,
		"unix:"+filepath.Join(paths.Base, "traffic-stats.socket"),
	)

	return output.String(), nil
}
