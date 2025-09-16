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

	contents := fmt.Sprintf(
		`
			user %s %s;
			load_module modules/ndk_http_module.so;
			load_module modules/ngx_http_js_module.so;
			load_module modules/ngx_http_lua_module.so;
			load_module modules/ngx_stream_module.so;
			worker_processes %d;
			pid %s/nginx.pid;
			error_log %s;
			
			events {
				worker_connections %d;
			}
			
			http {
				default_type %s;
				sendfile %s;
				server_tokens %s;
				client_max_body_size %dM;
				
				keepalive_timeout %d;
				proxy_connect_timeout %d;
				proxy_read_timeout %d;
				proxy_send_timeout %d;
				send_timeout %d;
				
				include %smime.types;
				%s
			}
			
			stream {
				%s
			}
		`,
		cfg.Nginx.RuntimeUser,
		cfg.Nginx.RuntimeUser,
		cfg.Nginx.WorkerProcesses,
		ctx.paths.AbsoluteConfig,
		p.getErrorLogPath(ctx.paths, logs),
		cfg.Nginx.WorkerConnections,
		cfg.Nginx.DefaultContentType,
		p.enabledFlag(cfg.Nginx.SendfileEnabled),
		p.enabledFlag(cfg.Nginx.ServerTokensEnabled),
		cfg.Nginx.MaximumBodySizeMb,
		cfg.Nginx.Timeouts.Keepalive,
		cfg.Nginx.Timeouts.Connect,
		cfg.Nginx.Timeouts.Read,
		cfg.Nginx.Timeouts.Send,
		cfg.Nginx.Timeouts.Send,
		ctx.paths,
		p.getHostIncludes(ctx.paths, ctx.hosts),
		p.getStreamIncludes(ctx.paths, ctx.streams),
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
		return fmt.Sprintf("%smain.log %s", paths.AbsoluteLogs, strings.ToLower(string(logs.ServerLogsLevel)))
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
		includes = append(includes, fmt.Sprintf("include %shost-%s.conf;", paths.AbsoluteConfig, h.ID))
	}

	return strings.Join(includes, "\n")
}

func (p *mainConfigurationFileProvider) getStreamIncludes(paths *Paths, streams []*stream.Stream) string {
	var includes []string
	for _, s := range streams {
		includes = append(includes, fmt.Sprintf("include %sstream-%s.conf;", paths.AbsoluteConfig, s.ID))
	}

	return strings.Join(includes, "\n")
}
