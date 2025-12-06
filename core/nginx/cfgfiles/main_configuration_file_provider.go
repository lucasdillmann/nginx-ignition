package cfgfiles

import (
	"fmt"
	"strings"

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
			moduleLines.WriteString("load_module modules/ngx_stream_module.so;\n")
		}

		streamLines.WriteString("stream {\n")
		streamLines.WriteString(p.getStreamIncludes(ctx.paths, ctx.streams))
		streamLines.WriteString("}\n")
	}

	if ctx.supportedFeatures.RunCodeType == DynamicSupportType {
		moduleLines.WriteString("load_module modules/ndk_http_module.so;\n")
		moduleLines.WriteString("load_module modules/ngx_http_js_module.so;\n")
		moduleLines.WriteString("load_module modules/ngx_http_lua_module.so;\n")
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
				default_type %s;
				sendfile %s;
				server_tokens %s;
				client_max_body_size %dM;
				tcp_nodelay %s;
				
				keepalive_timeout %ds;
				proxy_connect_timeout %ds;
				proxy_read_timeout %ds;
				proxy_send_timeout %ds;
				send_timeout %ds;
				client_body_timeout %ds;

				client_body_buffer_size %dk;
				client_header_buffer_size %dk;
				large_client_header_buffers %d %dk;
				output_buffers %d %dk;
				
				include %smime.types;
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
		cfg.Nginx.DefaultContentType,
		p.enabledFlag(cfg.Nginx.SendfileEnabled),
		p.enabledFlag(cfg.Nginx.ServerTokensEnabled),
		cfg.Nginx.MaximumBodySizeMb,
		p.enabledFlag(cfg.Nginx.TcpNoDelayEnabled),
		cfg.Nginx.Timeouts.Keepalive,
		cfg.Nginx.Timeouts.Connect,
		cfg.Nginx.Timeouts.Read,
		cfg.Nginx.Timeouts.Send,
		cfg.Nginx.Timeouts.Send,
		cfg.Nginx.Timeouts.ClientBody,
		cfg.Nginx.Buffers.ClientBodyKb,
		cfg.Nginx.Buffers.ClientHeaderKb,
		cfg.Nginx.Buffers.LargeClientHeader.Amount,
		cfg.Nginx.Buffers.LargeClientHeader.SizeKb,
		cfg.Nginx.Buffers.Output.Amount,
		cfg.Nginx.Buffers.Output.SizeKb,
		ctx.paths.Config,
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

func (p *mainConfigurationFileProvider) enabledFlag(value bool) string {
	if value {
		return "on"
	}

	return "off"
}

func (p *mainConfigurationFileProvider) getHostIncludes(paths *Paths, hosts []*host.Host) string {
	var includes []string
	for _, h := range hosts {
		includes = append(includes, fmt.Sprintf("include %shost-%s.conf;", paths.Config, h.ID))
	}

	return strings.Join(includes, "\n")
}

func (p *mainConfigurationFileProvider) getStreamIncludes(paths *Paths, streams []*stream.Stream) string {
	var includes []string

	for _, s := range streams {
		includes = append(includes, fmt.Sprintf("include %sstream-%s.conf;", paths.Config, s.ID))
	}

	return strings.Join(includes, "\n")
}
