package cfgfiles

import (
	"fmt"
	"strings"

	"dillmann.com.br/nginx-ignition/core/cache"
	"dillmann.com.br/nginx-ignition/core/common/log"
	"dillmann.com.br/nginx-ignition/core/common/ptr"
	"dillmann.com.br/nginx-ignition/core/host"
	"dillmann.com.br/nginx-ignition/core/settings"
	"dillmann.com.br/nginx-ignition/core/stream"
)

type mainConfigurationFileProvider struct{}

func newMainConfigurationFileProvider() *mainConfigurationFileProvider {
	return &mainConfigurationFileProvider{}
}

func (p *mainConfigurationFileProvider) provide(ctx *providerContext) ([]File, error) {
	cfg := ctx.settings

	moduleLines := strings.Builder{}
	if ctx.supportedFeatures.RunCodeType == DynamicSupportType {
		moduleLines.WriteString("load_module modules/ndk_http_module.so;\n")
		moduleLines.WriteString("load_module modules/ngx_http_js_module.so;\n")
		moduleLines.WriteString("load_module modules/ngx_http_lua_module.so;\n")
	}

	streamBlock := p.getStreamBlock(ctx, &moduleLines)
	mainHttpBlock := p.getMainHttpBlock(ctx, &moduleLines)

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
			
			%s
			%s
		`,
		cfg.Nginx.RuntimeUser,
		cfg.Nginx.RuntimeUser,
		moduleLines.String(),
		cfg.Nginx.WorkerProcesses,
		ctx.paths.Base,
		p.getErrorLogPath(ctx.paths, &cfg.Nginx.Logs),
		cfg.Nginx.WorkerConnections,
		mainHttpBlock,
		streamBlock,
	)

	return []File{
		{
			Name:     "nginx.conf",
			Contents: contents,
		},
	}, nil
}

func (p *mainConfigurationFileProvider) getStreamBlock(
	ctx *providerContext,
	modules *strings.Builder,
) string {
	if ctx.supportedFeatures.StreamType == NoneSupportType {
		return ""
	}

	if ctx.supportedFeatures.StreamType == DynamicSupportType {
		modules.WriteString("load_module modules/ngx_stream_module.so;\n")
	}

	streamIncludes := p.getStreamIncludes(ctx.paths, ctx.streams)
	if streamIncludes == "" {
		return ""
	}

	return fmt.Sprintf(
		`
			stream {
				%s
			}
		`,
		streamIncludes,
	)
}

func (p *mainConfigurationFileProvider) getStatusServerBlock(
	ctx *providerContext,
	modules *strings.Builder,
) string {
	cfg := ctx.settings.Nginx.API
	if !cfg.Enabled {
		return ""
	}

	switch ctx.supportedFeatures.API {
	case DynamicSupportType:
		modules.WriteString("load_module modules/ngx_http_api_module.so;\n")
	case NoneSupportType:
		log.Warnf("Nginx API cannot be enabled: API module is not available")
		return ""
	}

	return fmt.Sprintf(
		`
			server {
				listen %s:%d;
				location / {
					api write=%s;
				}
			}
		`,
		cfg.Address,
		cfg.Port,
		statusFlag(cfg.WriteEnabled),
	)
}

func (p *mainConfigurationFileProvider) getMainHttpBlock(
	ctx *providerContext,
	modules *strings.Builder,
) string {
	cfg := ctx.settings

	var customCfg string
	if cfg.Nginx.Custom != nil {
		customCfg = fmt.Sprintf("\n%s\n", *cfg.Nginx.Custom)
	}

	return fmt.Sprintf(
		`
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
				%s
			}
		`,
		statusFlag(cfg.Nginx.SendfileEnabled),
		statusFlag(cfg.Nginx.ServerTokensEnabled),
		statusFlag(cfg.Nginx.TcpNoDelayEnabled),
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
		p.getStatusServerBlock(ctx, modules),
	)
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
