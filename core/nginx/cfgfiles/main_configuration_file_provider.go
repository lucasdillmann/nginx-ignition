package cfgfiles

import (
	"fmt"
	"path/filepath"
	"strings"

	"dillmann.com.br/nginx-ignition/core/cache"
	"dillmann.com.br/nginx-ignition/core/common/ptr"
	"dillmann.com.br/nginx-ignition/core/common/runtime"
	"dillmann.com.br/nginx-ignition/core/host"
	"dillmann.com.br/nginx-ignition/core/settings"
	"dillmann.com.br/nginx-ignition/core/stream"
)

type mainConfigurationFileProvider struct {
	settingsCommands settings.Commands
}

func newMainConfigurationFileProvider(
	settingsCommands settings.Commands,
) *mainConfigurationFileProvider {
	return &mainConfigurationFileProvider{settingsCommands: settingsCommands}
}

func (p *mainConfigurationFileProvider) provide(ctx *providerContext) ([]File, error) {
	cfg, err := p.settingsCommands.Get(ctx.context)
	if err != nil {
		return nil, err
	}

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

	if ctx.supportedFeatures.StreamType == DynamicSupportType {
		_, _ = moduleLines.WriteString(
			"load_module modules/ngx_http_vhost_traffic_status_module.so;\n",
		)
	}

	var customCfg string
	if cfg.Nginx.Custom != nil {
		customCfg = fmt.Sprintf("\n%s\n", *cfg.Nginx.Custom)
	}

	userStatement := fmt.Sprintf("user %s %s;", cfg.Nginx.RuntimeUser, cfg.Nginx.RuntimeUser)
	if runtime.IsWindows() {
		userStatement = ""
	}

	contents := fmt.Sprintf(
		`
			%s
			%s
			worker_processes %d;
			pid "%snginx.pid";
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
				client_body_temp_path "%s" 1 2;
				proxy_temp_path "%s" 1 2;
				fastcgi_temp_path "%s" 1 2;
				scgi_temp_path "%s" 1 2;
				uwsgi_temp_path "%s" 1 2;

				default_type %s;
				include "%smime.types";
				%s
				%s
				%s
				%s
			}
			
			%s
		`,
		userStatement,
		moduleLines.String(),
		cfg.Nginx.WorkerProcesses,
		ctx.paths.Base,
		p.getErrorLogPath(ctx.paths, logs),
		cfg.Nginx.WorkerConnections,
		statusFlag(cfg.Nginx.SendfileEnabled),
		statusFlag(cfg.Nginx.ServerTokensEnabled),
		statusFlag(cfg.Nginx.TCPNoDelayEnabled),
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
		ctx.paths.Config,
		customCfg,
		p.getCacheDefinitions(ctx.paths, ctx.caches),
		p.getStatsDefinitions(ctx.paths, cfg.Nginx.Stats),
		p.getHostIncludes(ctx.paths, ctx.hosts),
		streamLines.String(),
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
		return fmt.Sprintf(
			"\"%smain.log\" %s",
			paths.Logs,
			strings.ToLower(string(logs.ServerLogsLevel)),
		)
	}

	return "off"
}

func (p *mainConfigurationFileProvider) getHostIncludes(paths *Paths, hosts []host.Host) string {
	includes := make([]string, 0, len(hosts))
	for _, h := range hosts {
		includes = append(includes, fmt.Sprintf("include \"%shost-%s.conf\";", paths.Config, h.ID))
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
			fmt.Sprintf("include \"%sstream-%s.conf\";", paths.Config, s.ID),
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
			storagePath = ptr.Of(paths.Cache + cacheIDNoDashes)
		}

		inactive := ""
		if c.InactiveSeconds != nil {
			inactive = fmt.Sprintf(" inactive=%ds", *c.InactiveSeconds)
		}

		maxSize := ""
		if c.MaximumSizeMB != nil {
			maxSize = fmt.Sprintf(" max_size=%dm", *c.MaximumSizeMB)
		}

		results = append(results, fmt.Sprintf(
			"proxy_cache_path \"%s\" levels=1:2 keys_zone=cache_%s:10m%s%s;",
			*storagePath,
			cacheIDNoDashes,
			inactive,
			maxSize,
		))
	}

	return strings.Join(results, "\n")
}

func (p *mainConfigurationFileProvider) getStatsDefinitions(
	paths *Paths,
	cfg *settings.NginxStatsSettings,
) string {
	if !cfg.Enabled {
		return ""
	}

	output := strings.Builder{}
	_, _ = fmt.Fprintf(
		&output,
		"\nvhost_traffic_status_zone shared:nginx-ignition-traffic-stats:%dm;\n",
		cfg.MaximumSizeMB,
	)

	if cfg.Persistent {
		dbLocation := cfg.DatabaseLocation
		if dbLocation == nil {
			dbLocation = ptr.Of(filepath.Join(paths.Base, "stats.db"))
		}

		_, _ = fmt.Fprintf(&output, "vhost_traffic_status_dump \"%s\" 5s;\n", *dbLocation)
	}

	_, _ = fmt.Fprintf(&output,
		`
		server { 
			root /dev/null;
            access_log off;
            listen unix:%s;
			
			location / {
				vhost_traffic_status_display;
				vhost_traffic_status_display_format json;
			}
        }
		`,
		filepath.Join(paths.Base, "traffic-stats.socket"),
	)

	return output.String()
}
