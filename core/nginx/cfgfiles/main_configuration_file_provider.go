package cfgfiles

import (
	"dillmann.com.br/nginx-ignition/core/host"
	"dillmann.com.br/nginx-ignition/core/settings"
	"fmt"
	"strings"
)

type mainConfigurationFileProvider struct {
	settingsRepository settings.Repository
}

func newMainConfigurationFileProvider(settingsRepository settings.Repository) *mainConfigurationFileProvider {
	return &mainConfigurationFileProvider{settingsRepository: settingsRepository}
}

func (p *mainConfigurationFileProvider) provide(basePath string, hosts []*host.Host) ([]output, error) {
	cfg, err := p.settingsRepository.Get()
	if err != nil {
		return nil, err
	}
	logs := cfg.Nginx.Logs

	contents := fmt.Sprintf(
		`
			load_module modules/ndk_http_module.so;
			load_module modules/ngx_http_js_module.so;
			load_module modules/ngx_http_lua_module.so;
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
				
				include %s/config/mime.types;
				%s
			}
		`,
		cfg.Nginx.WorkerProcesses,
		basePath,
		p.getErrorLogPath(basePath, logs),
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
		basePath,
		p.getHostIncludes(basePath, hosts),
	)

	return []output{
		{
			name:     "nginx.conf",
			contents: contents,
		},
	}, nil
}

func (p *mainConfigurationFileProvider) getErrorLogPath(basePath string, logs *settings.NginxLogsSettings) string {
	if logs.ServerLogsEnabled {
		return fmt.Sprintf("%s/logs/main.log %s", basePath, strings.ToLower(string(logs.ServerLogsLevel)))
	}

	return "off"
}

func (p *mainConfigurationFileProvider) enabledFlag(value bool) string {
	if value {
		return "on"
	}

	return "off"
}

func (p *mainConfigurationFileProvider) getHostIncludes(basePath string, hosts []*host.Host) string {
	var includes []string
	for _, h := range hosts {
		includes = append(includes, fmt.Sprintf("include %s/config/host-%s.conf;", basePath, h.ID))
	}

	return strings.Join(includes, "\n")
}
