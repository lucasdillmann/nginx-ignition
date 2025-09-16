package cfgfiles

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"dillmann.com.br/nginx-ignition/core/common/core_error"
	"dillmann.com.br/nginx-ignition/core/host"
	"dillmann.com.br/nginx-ignition/core/integration"
	"dillmann.com.br/nginx-ignition/core/settings"
)

type hostConfigurationFileProvider struct {
	integrationCommands *integration.Commands
	settingsRepository  settings.Repository
}

func newHostConfigurationFileProvider(
	settingsRepository settings.Repository,
	integrationCommands *integration.Commands,
) *hostConfigurationFileProvider {
	return &hostConfigurationFileProvider{
		integrationCommands: integrationCommands,
		settingsRepository:  settingsRepository,
	}
}

func (p *hostConfigurationFileProvider) provide(ctx *providerContext) ([]File, error) {
	var outputs []File
	for _, h := range ctx.hosts {
		if h.Enabled {
			output, err := p.buildHost(ctx.context, ctx.paths, h)
			if err != nil {
				return nil, err
			}

			outputs = append(outputs, *output)
		}
	}

	return outputs, nil
}

func (p *hostConfigurationFileProvider) buildHost(ctx context.Context, paths *Paths, h *host.Host) (*File, error) {
	var routes []string
	for _, r := range h.Routes {
		if r.Enabled {
			route, err := p.buildRoute(ctx, h, r, paths)
			if err != nil {
				return nil, err
			}

			routes = append(routes, route)
		}
	}

	serverNames, err := p.buildServerNames(h)
	if err != nil {
		return nil, err
	}

	httpsRedirect := ""
	if h.FeatureSet.RedirectHTTPToHTTPS {
		httpsRedirect = `if ($scheme = "http") { return 301 https://$server_name$request_uri; }`
	}

	http2 := ""
	if h.FeatureSet.HTTP2Support {
		http2 = "http2 on;"
	}

	bindings := h.Bindings
	if h.UseGlobalBindings {
		cfg, err := p.settingsRepository.Get(ctx)
		if err != nil {
			return nil, err
		}

		bindings = cfg.GlobalBindings
	}

	var contents []string
	for _, b := range bindings {
		binding, err := p.buildBinding(ctx, paths, h, b, routes, *serverNames, httpsRedirect, http2)
		if err != nil {
			return nil, err
		}
		contents = append(contents, binding)
	}

	return &File{
		Name:     fmt.Sprintf("host-%s.conf", h.ID),
		Contents: strings.Join(contents, "\n"),
	}, nil
}

func (p *hostConfigurationFileProvider) buildServerNames(h *host.Host) (*string, error) {
	serverNames := ""

	if h.DefaultServer {
		serverNames = "server_name _;"
	} else if len(h.DomainNames) > 0 {
		domainNames := make([]string, len(h.DomainNames))
		for index, domainName := range h.DomainNames {
			if domainName == nil {
				return nil, core_error.New("Unexpected null domain Name", false)
			}

			domainNames[index] = *domainName
		}

		serverNames = "server_name " + strings.Join(domainNames, " ") + ";"
	}

	return &serverNames, nil
}

func (p *hostConfigurationFileProvider) buildBinding(
	ctx context.Context,
	paths *Paths,
	h *host.Host,
	b *host.Binding,
	routes []string,
	serverNames, httpsRedirect, http2 string,
) (string, error) {
	listen := ""
	switch b.Type {
	case host.HttpBindingType:
		listen = fmt.Sprintf("listen %s:%d %s;", b.IP, b.Port, p.buildBindingAdditionalParams(h))
	case host.HttpsBindingType:
		listen = fmt.Sprintf(
			`
				listen %s:%d ssl %s;
				ssl_certificate %scertificate-%s.pem;
				ssl_certificate_key %scertificate-%s.pem;
				ssl_protocols TLSv1.2 TLSv1.3;
				ssl_ciphers HIGH:!aNULL:!MD5;
			`,
			b.IP,
			b.Port,
			p.buildBindingAdditionalParams(h),
			paths.Config,
			b.CertificateID,
			paths.Config,
			b.CertificateID,
		)
	}

	conditionalHttpsRedirect := ""
	if b.Type == host.HttpBindingType {
		conditionalHttpsRedirect = httpsRedirect
	}

	cfg, err := p.settingsRepository.Get(ctx)
	if err != nil {
		return "", err
	}

	logs := cfg.Nginx.Logs

	return fmt.Sprintf(
		`server {
			root /dev/null;
			access_log %s;
			error_log %s;
			gzip %s;
			client_max_body_size %dM;
			%s
			%s
			%s
			%s
			%s
			%s
		}`,
		p.flag(logs.AccessLogsEnabled, fmt.Sprintf("%shost-%s.access.log", paths.Logs, h.ID), "off"),
		p.flag(logs.ErrorLogsEnabled, fmt.Sprintf("%shost-%s.error.log %s", paths.Logs, h.ID, strings.ToLower(string(logs.ErrorLogsLevel))), "off"),
		p.flag(cfg.Nginx.GzipEnabled, "on", "off"),
		cfg.Nginx.MaximumBodySizeMb,
		p.flag(h.AccessListID != nil, fmt.Sprintf("include %saccess-list-%s.conf;", paths.Config, h.AccessListID), ""),
		conditionalHttpsRedirect,
		http2,
		listen,
		serverNames,
		strings.Join(routes, "\n"),
	), nil
}

func (p *hostConfigurationFileProvider) buildBindingAdditionalParams(h *host.Host) string {
	if h.DefaultServer {
		return "default_server"
	}

	return ""
}

func (p *hostConfigurationFileProvider) buildRoute(
	ctx context.Context,
	h *host.Host,
	r *host.Route,
	paths *Paths,
) (string, error) {
	switch r.Type {
	case host.StaticResponseRouteType:
		return p.buildStaticResponseRoute(r, h.FeatureSet, paths), nil
	case host.ProxyRouteType:
		return p.buildProxyRoute(r, h.FeatureSet, paths), nil
	case host.RedirectRouteType:
		return p.buildRedirectRoute(r, h.FeatureSet, paths), nil
	case host.IntegrationRouteType:
		return p.buildIntegrationRoute(ctx, r, h.FeatureSet, paths)
	case host.ExecuteCodeRouteType:
		return p.buildExecuteCodeRoute(h, r, paths), nil
	case host.StaticFilesRouteType:
		return p.buildStaticFilesRoute(r, paths), nil
	default:
		return "", fmt.Errorf("invalid route type: %s", r.Type)
	}
}

func (p *hostConfigurationFileProvider) buildStaticFilesRoute(r *host.Route, paths *Paths) string {
	normalizedSourcePath := r.SourcePath
	if !strings.HasSuffix(normalizedSourcePath, "/") {
		normalizedSourcePath += "/"
	}

	autoIndex := p.flag(r.Settings.DirectoryListingEnabled, "on", "off")

	return fmt.Sprintf(
		`location %s {
			rewrite  ^%s(.*) /$1 break;
			root %s;
			autoindex %s;
			autoindex_exact_size off;
			autoindex_format html;
			autoindex_localtime on;
			%s
		}`,
		normalizedSourcePath,
		normalizedSourcePath,
		*r.TargetURI,
		autoIndex,
		p.buildRouteSettings(r, paths),
	)
}

func (p *hostConfigurationFileProvider) buildStaticResponseRoute(
	r *host.Route,
	features host.FeatureSet,
	paths *Paths,
) string {
	var headers []string
	for key, value := range r.Response.Headers {
		headers = append(headers, fmt.Sprintf(`add_header "%s" "%s";`, key, value))
	}

	payload := ""
	if r.Response.Payload != nil && strings.TrimSpace(*r.Response.Payload) != "" {
		payload = strings.ReplaceAll(*r.Response.Payload, "\"", "\\\"")
		payload = fmt.Sprintf("\"%s\"", payload)
	}

	return fmt.Sprintf(
		`location %s {
			%s
			return %d %s;
			%s
			%s
		}`,
		r.SourcePath,
		strings.Join(headers, "\n"),
		r.Response.StatusCode,
		payload,
		p.buildRouteFeatures(features),
		p.buildRouteSettings(r, paths),
	)
}

func (p *hostConfigurationFileProvider) buildProxyRoute(
	r *host.Route,
	features host.FeatureSet,
	paths *Paths,
) string {
	return fmt.Sprintf(
		`location %s {
			%s
			%s
			%s
		}`,
		r.SourcePath,
		p.buildProxyPass(r),
		p.buildRouteFeatures(features),
		p.buildRouteSettings(r, paths),
	)
}

func (p *hostConfigurationFileProvider) buildIntegrationRoute(
	ctx context.Context,
	r *host.Route,
	features host.FeatureSet,
	paths *Paths,
) (string, error) {
	proxyUrl, err := p.integrationCommands.GetOptionUrlById(ctx, r.Integration.IntegrationID, r.Integration.OptionID)
	if err != nil {
		return "", err
	}

	if proxyUrl == nil {
		return "", core_error.New("Integration option not found", false)
	}

	return fmt.Sprintf(
		`location %s {
			%s
			%s
			%s
		}`,
		r.SourcePath,
		p.buildProxyPass(r, *proxyUrl),
		p.buildRouteFeatures(features),
		p.buildRouteSettings(r, paths),
	), nil
}

func (p *hostConfigurationFileProvider) buildRedirectRoute(
	r *host.Route,
	features host.FeatureSet,
	paths *Paths,
) string {
	return fmt.Sprintf(
		`location %s {
			return %d %s;
			%s
			%s
		}`,
		r.SourcePath,
		*r.RedirectCode,
		*r.TargetURI,
		p.buildRouteFeatures(features),
		p.buildRouteSettings(r, paths),
	)
}

func (p *hostConfigurationFileProvider) buildExecuteCodeRoute(h *host.Host, r *host.Route, paths *Paths) string {
	var headerBlock, routeBlock string
	switch r.SourceCode.Language {
	case host.JavascriptCodeLanguage:
		headerBlock = fmt.Sprintf("js_import route_%d from %shost-%s-route-%d.js;", r.Priority, paths.Config, h.ID, r.Priority)
		routeBlock = fmt.Sprintf("js_content route_%d.%s;", r.Priority, *r.SourceCode.MainFunction)
	case host.LuaCodeLanguage:
		routeBlock = fmt.Sprintf(
			`content_by_lua_block {
				%s
			}`,
			r.SourceCode.Contents,
		)
	}

	return fmt.Sprintf(
		`%s
		location %s {
			%s
			%s
			%s
		}`,
		headerBlock,
		r.SourcePath,
		routeBlock,
		p.buildRouteFeatures(h.FeatureSet),
		p.buildRouteSettings(r, paths),
	)
}

func (p *hostConfigurationFileProvider) buildRouteFeatures(features host.FeatureSet) string {
	if features.WebsocketSupport {
		return `
			proxy_http_version 1.1;
			proxy_set_header Upgrade $http_upgrade;
			proxy_set_header Connection "upgrade";
		`
	}

	return ""
}

func (p *hostConfigurationFileProvider) buildProxyPass(r *host.Route, uri ...string) string {
	targetUri := r.TargetURI
	if len(uri) > 0 {
		targetUri = &uri[0]
	}

	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("proxy_pass %s;", *targetUri))

	if r.Settings.KeepOriginalDomainName {
		u, _ := url.Parse(*targetUri)
		builder.WriteString(fmt.Sprintf("\nproxy_set_header Host %s;", u.Host))
	}

	return builder.String()
}

func (p *hostConfigurationFileProvider) buildRouteSettings(r *host.Route, paths *Paths) string {
	builder := strings.Builder{}
	if r.Settings.ProxySSLServerName {
		builder.WriteString("proxy_ssl_server_name on;")
	}

	if r.Settings.IncludeForwardHeaders {
		builder.WriteString(`
			proxy_set_header x-forwarded-for $proxy_add_x_forwarded_for;
			proxy_set_header x-forwarded-host $host;
			proxy_set_header x-forwarded-proto $scheme;
			proxy_set_header x-forwarded-scheme $scheme;
			proxy_set_header x-forwarded-port $server_port;
			proxy_set_header x-real-ip $remote_addr;
		`)
	}

	if r.Settings.Custom != nil {
		builder.WriteString("\n")
		builder.WriteString(*r.Settings.Custom)
	}

	if r.AccessListID != nil {
		builder.WriteString(fmt.Sprintf("\ninclude %saccess-list-%s.conf;", paths.Config, *r.AccessListID))
	}

	return builder.String()
}

func (p *hostConfigurationFileProvider) flag(enabled bool, trueValue, falseValue string) string {
	if enabled {
		return trueValue
	}

	return falseValue
}
