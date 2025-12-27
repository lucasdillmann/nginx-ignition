package cfgfiles

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/binding"
	"dillmann.com.br/nginx-ignition/core/cache"
	"dillmann.com.br/nginx-ignition/core/common/coreerror"
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
	outputs := make([]File, 0)
	for _, h := range ctx.hosts {
		if h.Enabled {
			output, err := p.buildHost(ctx, &h)
			if err != nil {
				return nil, err
			}

			outputs = append(outputs, *output)
		}
	}

	return outputs, nil
}

func (p *hostConfigurationFileProvider) buildHost(ctx *providerContext, h *host.Host) (*File, error) {
	routes := make([]string, 0)
	for _, r := range h.Routes {
		if r.Enabled {
			route, err := p.buildRoute(ctx, h, &r)
			if err != nil {
				return nil, err
			}

			routes = append(routes, route)
		}
	}

	serverNames := p.buildServerNames(h)

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
		cfg, err := p.settingsRepository.Get(ctx.context)
		if err != nil {
			return nil, err
		}

		bindings = cfg.GlobalBindings
	}

	contents := make([]string, 0)
	for _, b := range bindings {
		binding, err := p.buildBinding(ctx, h, &b, routes, serverNames, httpsRedirect, http2)
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

func (p *hostConfigurationFileProvider) buildServerNames(h *host.Host) string {
	if h.DefaultServer {
		return "server_name _;"
	}

	return "server_name " + strings.Join(h.DomainNames, " ") + ";"
}

func (p *hostConfigurationFileProvider) buildBinding(
	ctx *providerContext,
	h *host.Host,
	b *binding.Binding,
	routes []string,
	serverNames, httpsRedirect, http2 string,
) (string, error) {
	listen := ""
	switch b.Type {
	case binding.HttpBindingType:
		listen = fmt.Sprintf("listen %s:%d %s;", b.IP, b.Port, p.buildBindingAdditionalParams(h))
	case binding.HttpsBindingType:
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
			ctx.paths.Config,
			b.CertificateID,
			ctx.paths.Config,
			b.CertificateID,
		)
	}

	conditionalHttpsRedirect := ""
	if b.Type == binding.HttpBindingType {
		conditionalHttpsRedirect = httpsRedirect
	}

	cfg, err := p.settingsRepository.Get(ctx.context)
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
			%s
		}`,
		flag(logs.AccessLogsEnabled, fmt.Sprintf("%shost-%s.access.log", ctx.paths.Logs, h.ID), "off"),
		flag(logs.ErrorLogsEnabled, fmt.Sprintf("%shost-%s.error.log %s", ctx.paths.Logs, h.ID, strings.ToLower(string(logs.ErrorLogsLevel))), "off"),
		statusFlag(cfg.Nginx.GzipEnabled),
		cfg.Nginx.MaximumBodySizeMb,
		flag(h.AccessListID != nil, fmt.Sprintf("include %saccess-list-%s.conf;", ctx.paths.Config, h.AccessListID), ""),
		p.buildCacheConfig(ctx.caches, h.CacheID),
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
	ctx *providerContext,
	h *host.Host,
	r *host.Route,
) (string, error) {
	switch r.Type {
	case host.StaticResponseRouteType:
		return p.buildStaticResponseRoute(ctx, h, r), nil
	case host.ProxyRouteType:
		return p.buildProxyRoute(ctx, r, h.FeatureSet), nil
	case host.RedirectRouteType:
		return p.buildRedirectRoute(ctx, r, h.FeatureSet), nil
	case host.IntegrationRouteType:
		return p.buildIntegrationRoute(ctx, r, h.FeatureSet)
	case host.ExecuteCodeRouteType:
		return p.buildExecuteCodeRoute(ctx, h, r), nil
	case host.StaticFilesRouteType:
		return p.buildStaticFilesRoute(ctx, r), nil
	default:
		return "", fmt.Errorf("invalid route type: %s", r.Type)
	}
}

func (p *hostConfigurationFileProvider) buildStaticFilesRoute(ctx *providerContext, r *host.Route) string {
	normalizedSourcePath := r.SourcePath
	if !strings.HasSuffix(normalizedSourcePath, "/") {
		normalizedSourcePath += "/"
	}

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
		statusFlag(r.Settings.DirectoryListingEnabled),
		p.buildRouteSettings(ctx, r),
	)
}

func (p *hostConfigurationFileProvider) buildStaticResponseRoute(
	ctx *providerContext,
	h *host.Host,
	r *host.Route,
) string {
	headers := ""
	payloadFilePath := fmt.Sprintf("/host-%s-route-%d.payload", h.ID, r.Priority)

	for key, value := range r.Response.Headers {
		headers += fmt.Sprintf(`add_header "%s" "%s" always;`, key, value) + "\n"
	}

	return fmt.Sprintf(
		`
		location @route_%d/static_payload {
			internal;
			%s
			root %s;
			try_files %s =%d;
		}

		location %s {
			%s
			error_page 599 =%d @route_%d/static_payload;
			%s
			%s
			return 599;
		}`,
		r.Priority,
		headers,
		ctx.paths.Config,
		payloadFilePath,
		r.Response.StatusCode,
		r.SourcePath,
		headers,
		r.Response.StatusCode,
		r.Priority,
		p.buildRouteFeatures(h.FeatureSet),
		p.buildRouteSettings(ctx, r),
	)
}

func (p *hostConfigurationFileProvider) buildProxyRoute(
	ctx *providerContext,
	r *host.Route,
	features host.FeatureSet,
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
		p.buildRouteSettings(ctx, r),
	)
}

func (p *hostConfigurationFileProvider) buildIntegrationRoute(
	ctx *providerContext,
	r *host.Route,
	features host.FeatureSet,
) (string, error) {
	proxyUrl, dnsResolvers, err := p.integrationCommands.GetOptionURL(ctx.context, r.Integration.IntegrationID, r.Integration.OptionID)
	if err != nil {
		return "", err
	}

	if proxyUrl == nil {
		msg := fmt.Sprintf("Integration option not found: %s", r.Integration.OptionID)
		return "", coreerror.New(msg, false)
	}

	dnsConfig := ""
	if len(dnsResolvers) > 0 {
		ips := strings.Join(dnsResolvers, " ")
		dnsConfig = fmt.Sprintf("resolver %s valid=5s;", ips)
	}

	return fmt.Sprintf(
		`location %s {
			%s
			%s
			%s
			%s
		}`,
		r.SourcePath,
		dnsConfig,
		p.buildProxyPass(r, *proxyUrl),
		p.buildRouteFeatures(features),
		p.buildRouteSettings(ctx, r),
	), nil
}

func (p *hostConfigurationFileProvider) buildRedirectRoute(
	ctx *providerContext,
	r *host.Route,
	features host.FeatureSet,
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
		p.buildRouteSettings(ctx, r),
	)
}

func (p *hostConfigurationFileProvider) buildExecuteCodeRoute(ctx *providerContext, h *host.Host, r *host.Route) string {
	var headerBlock, routeBlock string
	switch r.SourceCode.Language {
	case host.JavascriptCodeLanguage:
		headerBlock = fmt.Sprintf("js_import route_%d from %shost-%s-route-%d.js;", r.Priority, ctx.paths.Config, h.ID, r.Priority)
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
		p.buildRouteSettings(ctx, r),
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

	if targetUri == nil {
		return ""
	}

	builder := strings.Builder{}
	fmt.Fprintf(&builder, "proxy_pass %s;", *targetUri)

	if r.Settings.KeepOriginalDomainName {
		u, _ := url.Parse(*targetUri)
		fmt.Fprintf(&builder, "\nproxy_set_header Host %s;", u.Host)
	}

	return builder.String()
}

func (p *hostConfigurationFileProvider) buildRouteSettings(ctx *providerContext, r *host.Route) string {
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
		fmt.Fprintf(&builder, "\ninclude %saccess-list-%s.conf;", ctx.paths.Config, *r.AccessListID)
	}

	builder.WriteString(p.buildCacheConfig(ctx.caches, r.CacheID))

	return builder.String()
}

func (p *hostConfigurationFileProvider) buildCacheConfig(caches []cache.Cache, cacheID *uuid.UUID) string {
	if cacheID == nil || len(caches) == 0 {
		return ""
	}

	var c *cache.Cache
	for _, item := range caches {
		if item.ID == *cacheID {
			c = &item
			break
		}
	}

	if c == nil {
		return ""
	}

	builder := strings.Builder{}
	builder.WriteString("\n")

	cacheIDNoDashes := strings.ReplaceAll(c.ID.String(), "-", "")
	fmt.Fprintf(&builder, "proxy_cache cache_%s;", cacheIDNoDashes)

	p.appendCacheDurations(&builder, c)
	p.appendCacheMethods(&builder, c)
	p.appendCacheStandardOptions(&builder, c)
	p.appendCacheLock(&builder, c)
	p.appendCacheBypassRules(&builder, c)
	p.appendCacheFileExtensions(&builder, c)

	return builder.String()
}

func (p *hostConfigurationFileProvider) appendCacheDurations(builder *strings.Builder, c *cache.Cache) {
	for _, d := range c.Durations {
		fmt.Fprintf(
			builder,
			"\nproxy_cache_valid %s %ds;",
			strings.Join(d.StatusCodes, " "),
			d.ValidTimeSeconds,
		)
	}
}

func (p *hostConfigurationFileProvider) appendCacheMethods(builder *strings.Builder, c *cache.Cache) {
	if len(c.AllowedMethods) > 0 {
		methods := make([]string, len(c.AllowedMethods))
		for index, method := range c.AllowedMethods {
			methods[index] = strings.ToLower(string(method))
		}

		fmt.Fprintf(builder, "\nproxy_cache_methods %s;", strings.Join(methods, " "))
	}
}

func (p *hostConfigurationFileProvider) appendCacheStandardOptions(builder *strings.Builder, c *cache.Cache) {
	fmt.Fprintf(builder, "\nproxy_cache_min_uses %d;", c.MinimumUsesBeforeCaching)
	fmt.Fprintf(builder, "\nproxy_cache_background_update %s;", statusFlag(c.BackgroundUpdate))
	fmt.Fprintf(builder, "\nproxy_cache_revalidate %s;", statusFlag(c.Revalidate))

	staleConfig := offFlag

	if len(c.UseStale) > 0 {
		staleOptions := make([]string, len(c.UseStale))
		for index, option := range c.UseStale {
			staleOptions[index] = strings.ToLower(string(option))
		}

		staleConfig = strings.Join(staleOptions, " ")
	}

	fmt.Fprintf(builder, "\nproxy_cache_use_stale %s;", staleConfig)
}

func (p *hostConfigurationFileProvider) appendCacheLock(builder *strings.Builder, c *cache.Cache) {
	if c.ConcurrencyLock.Enabled {
		builder.WriteString("\nproxy_cache_lock on;")
		if c.ConcurrencyLock.TimeoutSeconds != nil {
			fmt.Fprintf(builder, "\nproxy_cache_lock_timeout %ds;", *c.ConcurrencyLock.TimeoutSeconds)
		}
		if c.ConcurrencyLock.AgeSeconds != nil {
			fmt.Fprintf(builder, "\nproxy_cache_lock_age %ds;", *c.ConcurrencyLock.AgeSeconds)
		}
	}
}

func (p *hostConfigurationFileProvider) appendCacheBypassRules(builder *strings.Builder, c *cache.Cache) {
	for _, rule := range c.BypassRules {
		fmt.Fprintf(builder, "\nproxy_cache_bypass %s;", rule)
	}

	for _, rule := range c.NoCacheRules {
		fmt.Fprintf(builder, "\nproxy_no_cache %s;", rule)
	}
}

func (p *hostConfigurationFileProvider) appendCacheFileExtensions(builder *strings.Builder, c *cache.Cache) {
	if len(c.FileExtensions) == 0 {
		return
	}

	extensions := make([]string, len(c.FileExtensions))
	for index, extension := range c.FileExtensions {
		extensions[index] = strings.ReplaceAll(
			strings.TrimSpace(extension),
			".", "\\.",
		)
	}

	builder.WriteString("\n")
	fmt.Fprintf(
		builder,
		`
			if ($uri !~* "%s") { 
				set $__no_cache_allowed_extension 1; 
			}
			proxy_no_cache $__no_cache_allowed_extension;
		`,
		fmt.Sprintf("\\.(%s)$", strings.Join(extensions, "|")),
	)
}
