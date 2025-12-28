package cfgfiles

import (
	"fmt"
	"strings"

	"dillmann.com.br/nginx-ignition/core/cache"
	"dillmann.com.br/nginx-ignition/core/common/ptr"
	"dillmann.com.br/nginx-ignition/core/host"
	"dillmann.com.br/nginx-ignition/core/settings"
	"dillmann.com.br/nginx-ignition/core/stream"
)

type mainConfigurationFileProvider struct {
	settingsRepository settings.Repository
}

func newMainConfigurationFileProvider(settingsRepository settings.Repository) *mainConfigurationFileProvider {
	return &mainConfigurationFileProvider{settingsRepository: settingsRepository}
}

func (p *mainConfigurationFileProvider) provide(ctx *providerContext) ([]File, error) {
	cfg, err := p.settingsRepository.Get(ctx.context)
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

	var customCfg string
	if cfg.Nginx.Custom != nil {
		customCfg = fmt.Sprintf("\n%s\n", *cfg.Nginx.Custom)
	}

	contents := fmt.Sprintf(
		`
			user %s %s;
			%s
			worker_processes %d;
			pid %snginx.pid;
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

				default_type %s;
				include %smime.types;
				%s
				%s
				%s
			}
			
			%s
		`,
		cfg.Nginx.RuntimeUser,
		cfg.Nginx.RuntimeUser,
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
		cfg.Nginx.DefaultContentType,
		ctx.paths.Config,
		customCfg,
		p.getCacheDefinitions(ctx.paths, ctx.caches),
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

func (p *mainConfigurationFileProvider) getErrorLogPath(paths *Paths, logs *settings.NginxLogsSettings) string {
	if logs.ServerLogsEnabled {
		return fmt.Sprintf("%smain.log %s", paths.Logs, strings.ToLower(string(logs.ServerLogsLevel)))
	}

	return "off"
}

func (p *mainConfigurationFileProvider) getHostIncludes(paths *Paths, hosts []host.Host) string {
	includes := make([]string, 0)
	for _, h := range hosts {
		includes = append(includes, fmt.Sprintf("include %shost-%s.conf;", paths.Config, h.ID))
	}

	return strings.Join(includes, "\n")
}

func (p *mainConfigurationFileProvider) getStreamIncludes(paths *Paths, streams []stream.Stream) string {
	includes := make([]string, 0)

	for _, s := range streams {
		includes = append(includes, fmt.Sprintf("include %sstream-%s.conf;", paths.Config, s.ID))
	}

	return strings.Join(includes, "\n")
}

func (p *mainConfigurationFileProvider) getCacheDefinitions(paths *Paths, caches []cache.Cache) string {
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
			"proxy_cache_path %s levels=1:2 keys_zone=cache_%s:10m%s%s;",
			*storagePath,
			cacheIDNoDashes,
			inactive,
			maxSize,
		))
	}

	return strings.Join(results, "\n")
}
