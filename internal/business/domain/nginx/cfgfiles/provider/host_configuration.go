package provider

import (
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"

	"github.com/google/uuid"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/coreerror"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/i18n"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/binding"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/cache"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/host"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/integration"
)

const (
	httpScheme  = "http://"
	httpsScheme = "https://"
)

type hostConfigurationProvider struct {
	integrationCommands integration.Commands
}

func newHostConfigurationProvider(
	integrationCommands integration.Commands,
) *hostConfigurationProvider {
	return &hostConfigurationProvider{
		integrationCommands: integrationCommands,
	}
}

func (p *hostConfigurationProvider) Provide(ctx *Context) ([]File, error) {
	outputs := make([]File, 0)
	for _, h := range ctx.Hosts {
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

func (p *hostConfigurationProvider) buildHost(
	ctx *Context,
	h *host.Host,
) (*File, error) {
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
		bindings = ctx.Cfg.GlobalBindings
	}

	stats := ""
	statsCfg := ctx.Cfg.Nginx.Stats

	if statsCfg.Enabled {
		stats = nginxSprintf(
			`
			set $stats_host_id %s;
			vhost_traffic_status %s;
			vhost_traffic_status_filter_by_set_key $stats_host_id hosts;
			vhost_traffic_status_filter_by_set_key $geoip_country_code countryCode@host:$stats_host_id;
			vhost_traffic_status_filter_by_set_key $geoip_city_name city@host:$stats_host_id;
			vhost_traffic_status_filter_by_set_key $stats_user_agent userAgent@host:$stats_host_id;
			vhost_traffic_status_filter_by_set_key $geoip_country_code countryCode@domain:$server_name;
			vhost_traffic_status_filter_by_set_key $geoip_city_name city@domain:$server_name;
			vhost_traffic_status_filter_by_set_key $stats_user_agent userAgent@domain:$server_name;
			`,
			h.ID,
			directiveFragment(statusFlag(statsCfg.AllHosts || h.FeatureSet.StatsEnabled)),
		)
	}

	contents := make([]string, 0)
	for _, b := range bindings {
		b, err := p.buildBinding(ctx, h, &b, routes, serverNames, httpsRedirect, http2, stats)
		if err != nil {
			return nil, err
		}
		contents = append(contents, b)
	}

	return &File{
		Name:     fmt.Sprintf("host-%s.conf", h.ID),
		Contents: strings.Join(contents, "\n"),
	}, nil
}

func (p *hostConfigurationProvider) buildServerNames(h *host.Host) string {
	if h.DefaultServer {
		return "server_name _;"
	}

	return nginxSprintfArgs("server_name %s;", h.DomainNames...)
}

func (p *hostConfigurationProvider) buildBinding(
	ctx *Context,
	h *host.Host,
	b *binding.Binding,
	routes []string,
	serverNames, httpsRedirect, http2, stats string,
) (string, error) {
	bindingAddress := net.JoinHostPort(b.IP, strconv.Itoa(b.Port))
	certificateID := ""
	if b.CertificateID != nil {
		certificateID = b.CertificateID.String()
	}

	listen := ""
	switch b.Type {
	case binding.HTTPBindingType:
		listen = nginxSprintf(
			"listen %s %s;",
			bindingAddress,
			directiveFragment(p.buildBindingAdditionalParams(h)),
		)
	case binding.HTTPSBindingType:
		listen = nginxSprintf(
			`
				listen %s ssl %s;
				ssl_certificate %s;
				ssl_certificate_key %s;
				ssl_protocols TLSv1.2 TLSv1.3;
				ssl_ciphers HIGH:!aNULL:!MD5;
			`,
			bindingAddress,
			directiveFragment(p.buildBindingAdditionalParams(h)),
			ctx.Paths.Config+"certificate-"+certificateID+".pem",
			ctx.Paths.Config+"certificate-"+certificateID+".pem",
		)
	default:
		return "", fmt.Errorf("invalid binding type: %s", b.Type)
	}

	conditionalHTTPSRedirect := ""
	if b.Type == binding.HTTPBindingType {
		conditionalHTTPSRedirect = httpsRedirect
	}

	logs := ctx.Cfg.Nginx.Logs
	accessLog := "access_log off;"
	if logs.AccessLogsEnabled {
		accessLog = nginxSprintf(
			"access_log %s;",
			ctx.Paths.Logs+"host-"+h.ID.String()+".access.log",
		)
	}

	errorLog := "error_log off;"
	if logs.ErrorLogsEnabled {
		errorLog = nginxSprintf(
			"error_log %s %s;",
			ctx.Paths.Logs+"host-"+h.ID.String()+".error.log",
			strings.ToLower(string(logs.ErrorLogsLevel)),
		)
	}

	accessList := ""
	if h.AccessListID != nil {
		accessList = nginxSprintf(
			"include %s;",
			ctx.Paths.Config+"access-list-"+h.AccessListID.String()+".conf",
		)
	}

	return nginxSprintf(
		`server {
			root /dev/null;
			%s
			%s
			gzip %s;
			client_max_body_size %dM;
			%s
			%s
			%s
			%s
			%s
			%s
			%s
			%s
		}`,
		directiveFragment(accessLog),
		directiveFragment(errorLog),
		directiveFragment(statusFlag(ctx.Cfg.Nginx.GzipEnabled)),
		ctx.Cfg.Nginx.MaximumBodySizeMb,
		directiveFragment(accessList),
		directiveFragment(p.buildCacheConfig(ctx.Caches, h.CacheID)),
		directiveFragment(conditionalHTTPSRedirect),
		directiveFragment(http2),
		directiveFragment(stats),
		directiveFragment(listen),
		directiveFragment(serverNames),
		directiveFragment(strings.Join(routes, "\n")),
	), nil
}

func (p *hostConfigurationProvider) buildBindingAdditionalParams(h *host.Host) string {
	if h.DefaultServer {
		return "default_server"
	}

	return ""
}

func (p *hostConfigurationProvider) buildRoute(
	ctx *Context,
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
		return p.buildExecuteCodeRoute(ctx, h, r)
	case host.StaticFilesRouteType:
		return p.buildStaticFilesRoute(ctx, r), nil
	default:
		return "", fmt.Errorf("invalid route type: %s", r.Type)
	}
}

func (p *hostConfigurationProvider) buildStaticFilesRoute(
	ctx *Context,
	r *host.Route,
) string {
	normalizedSourcePath := r.SourcePath
	if !strings.HasSuffix(normalizedSourcePath, "/") {
		normalizedSourcePath += "/"
	}

	indexFile := ""
	if r.Settings.IndexFile != nil && strings.TrimSpace(*r.Settings.IndexFile) != "" {
		indexFile = nginxSprintf("index %s;", *r.Settings.IndexFile)
	}

	rewritePattern := "^" + normalizedSourcePath + "(.*)$"

	return nginxSprintf(
		`location %s {
			rewrite %s %s break;
			root %s;
			%s
			autoindex %s;
			autoindex_exact_size off;
			autoindex_format html;
			autoindex_localtime on;
			%s
		}`,
		normalizedSourcePath,
		rewritePattern,
		"/$1",
		*r.TargetURI,
		directiveFragment(indexFile),
		directiveFragment(statusFlag(r.Settings.DirectoryListingEnabled)),
		directiveFragment(p.buildRouteSettings(ctx, r)),
	)
}

func (p *hostConfigurationProvider) buildStaticResponseRoute(
	ctx *Context,
	h *host.Host,
	r *host.Route,
) string {
	headers := ""
	payloadFilePath := fmt.Sprintf("/host-%s-route-%d.payload", h.ID, r.Priority)

	for key, value := range r.Response.Headers {
		headers += nginxSprintf(
			"add_header %s %s always;\n",
			key,
			sanitizeHeaderValue(value),
		)
	}

	return nginxSprintf(
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
		directiveFragment(headers),
		ctx.Paths.Config,
		payloadFilePath,
		r.Response.StatusCode,
		r.SourcePath,
		directiveFragment(headers),
		r.Response.StatusCode,
		r.Priority,
		directiveFragment(p.buildRouteFeatures(h.FeatureSet, r.Protocol)),
		directiveFragment(p.buildRouteSettings(ctx, r)),
	)
}

func (p *hostConfigurationProvider) buildProxyRoute(
	ctx *Context,
	r *host.Route,
	features host.FeatureSet,
) string {
	return nginxSprintf(
		`location %s {
			%s
			%s
			%s
			%s
		}`,
		r.SourcePath,
		directiveFragment(p.buildProxyPass(r)),
		directiveFragment(p.buildProtocolProxyVersion(r)),
		directiveFragment(p.buildRouteFeatures(features, r.Protocol)),
		directiveFragment(p.buildRouteSettings(ctx, r)),
	)
}

func (p *hostConfigurationProvider) buildIntegrationRoute(
	ctx *Context,
	r *host.Route,
	features host.FeatureSet,
) (string, error) {
	proxyURL, dnsResolvers, err := p.integrationCommands.GetOptionURL(
		ctx.Context,
		r.Integration.IntegrationID,
		r.Integration.OptionID,
	)
	if err != nil {
		return "", err
	}

	if proxyURL == nil {
		return "", coreerror.New(
			i18n.M(ctx.Context, i18n.K.CoreNginxCfgfilesOptionNotFound).
				V("optionId", r.Integration.OptionID),
			true,
		)
	}

	if r.Integration.UseHTTPS {
		proxyURL = new(strings.Replace(*proxyURL, httpScheme, httpsScheme, 1))
	}

	if r.TargetURI != nil && strings.TrimSpace(*r.TargetURI) != "" {
		proxyURL = new(*proxyURL + *r.TargetURI)
	}

	dnsConfig := ""
	if len(dnsResolvers) > 0 {
		dnsConfig = nginxSprintfArgs("resolver %s valid=5s;", dnsResolvers...)
	}

	return nginxSprintf(
		`location %s {
			%s
			%s
			%s
			%s
			%s
		}`,
		r.SourcePath,
		directiveFragment(dnsConfig),
		directiveFragment(p.buildProxyPass(r, *proxyURL)),
		directiveFragment(p.buildProtocolProxyVersion(r)),
		directiveFragment(p.buildRouteFeatures(features, r.Protocol)),
		directiveFragment(p.buildRouteSettings(ctx, r)),
	), nil
}

func (p *hostConfigurationProvider) buildRedirectRoute(
	ctx *Context,
	r *host.Route,
	features host.FeatureSet,
) string {
	return nginxSprintf(
		`location %s {
			return %d %s;
			%s
			%s
		}`,
		r.SourcePath,
		*r.RedirectCode,
		*r.TargetURI,
		directiveFragment(p.buildRouteFeatures(features, r.Protocol)),
		directiveFragment(p.buildRouteSettings(ctx, r)),
	)
}

func (p *hostConfigurationProvider) buildExecuteCodeRoute(
	ctx *Context,
	h *host.Host,
	r *host.Route,
) (string, error) {
	var headerBlock, routeBlock string
	switch r.SourceCode.Language {
	case host.JavascriptCodeLanguage:
		headerBlock = nginxSprintf(
			"js_import route_%d from %s;",
			r.Priority,
			ctx.Paths.Config+"host-"+h.ID.String()+"-route-"+strconv.Itoa(r.Priority)+".js",
		)
		routeBlock = nginxSprintf(
			"js_content %s;",
			"route_"+strconv.Itoa(r.Priority)+"."+*r.SourceCode.MainFunction,
		)
	case host.LuaCodeLanguage:
		routeBlock = nginxSprintf(
			`content_by_lua_block {
				%s
			}`,
			rawConfigFragment(r.SourceCode.Contents),
		)
	default:
		return "", fmt.Errorf("invalid language: %s", r.SourceCode.Language)
	}

	return nginxSprintf(
		`%s
		location %s {
			%s
			%s
			%s
		}`,
		directiveFragment(headerBlock),
		r.SourcePath,
		directiveFragment(routeBlock),
		directiveFragment(p.buildRouteFeatures(h.FeatureSet, r.Protocol)),
		directiveFragment(p.buildRouteSettings(ctx, r)),
	), nil
}

func (p *hostConfigurationProvider) buildRouteFeatures(
	features host.FeatureSet,
	protocol host.RouteProtocol,
) string {
	if features.WebsocketSupport && protocol == host.HTTP11RouteProtocol {
		return `
			proxy_set_header Upgrade $http_upgrade;
			proxy_set_header Connection "upgrade";
		`
	}

	return ""
}

func (p *hostConfigurationProvider) buildProxyPass(r *host.Route, uri ...string) string {
	targetURI := r.TargetURI
	if len(uri) > 0 {
		targetURI = &uri[0]
	}

	if targetURI == nil {
		return ""
	}

	builder := strings.Builder{}
	grpcProtocol := r.Protocol == host.GRPCRouteProtocol

	if grpcProtocol {
		nginxFprintf(&builder, "grpc_pass %s;", p.toGrpcURL(*targetURI))
	} else {
		nginxFprintf(&builder, "proxy_pass %s;", *targetURI)
	}

	if r.Settings.KeepOriginalDomainName {
		u, _ := url.Parse(*targetURI)
		if grpcProtocol {
			nginxFprintf(&builder, "\ngrpc_set_header Host %s;", u.Host)
		} else {
			nginxFprintf(&builder, "\nproxy_set_header Host %s;", u.Host)
		}
	}

	return builder.String()
}

func (p *hostConfigurationProvider) buildProtocolProxyVersion(r *host.Route) string {
	switch r.Protocol {
	case host.GRPCRouteProtocol:
		return ""
	case host.HTTP10RouteProtocol:
		return "proxy_http_version 1.0;"
	default:
		return "proxy_http_version 1.1;"
	}
}

func (p *hostConfigurationProvider) toGrpcURL(uri string) string {
	switch {
	case strings.HasPrefix(uri, httpsScheme):
		return "grpcs://" + strings.TrimPrefix(uri, httpsScheme)
	case strings.HasPrefix(uri, httpScheme):
		return "grpc://" + strings.TrimPrefix(uri, httpScheme)
	default:
		return uri
	}
}

func (p *hostConfigurationProvider) buildRouteSettings(
	ctx *Context,
	r *host.Route,
) string {
	builder := strings.Builder{}

	if r.Settings.ProxySSLServerName {
		_, _ = builder.WriteString("proxy_ssl_server_name on;\n")
	}

	if r.Settings.IgnoreSSLErrors {
		_, _ = builder.WriteString("proxy_ssl_verify off;\n")
	}

	if r.Settings.IncludeForwardHeaders {
		prefix := "proxy_set_header "
		if r.Protocol == host.GRPCRouteProtocol {
			prefix = "grpc_set_header "
		}

		_, _ = builder.WriteString(prefix + "x-forwarded-for $proxy_add_x_forwarded_for;\n")
		_, _ = builder.WriteString(prefix + "x-forwarded-host $host;\n")
		_, _ = builder.WriteString(prefix + "x-forwarded-proto $scheme;\n")
		_, _ = builder.WriteString(prefix + "x-forwarded-scheme $scheme;\n")
		_, _ = builder.WriteString(prefix + "x-forwarded-port $server_port;\n")
		_, _ = builder.WriteString(prefix + "x-real-ip $remote_addr;\n")
	}

	if r.Settings.Custom != nil {
		_, _ = builder.WriteString("\n")
		_, _ = builder.WriteString(string(rawConfigFragment(*r.Settings.Custom)))
	}

	if r.AccessListID != nil {
		nginxFprintf(
			&builder,
			"\ninclude %s;",
			ctx.Paths.Config+"access-list-"+r.AccessListID.String()+".conf",
		)
	}

	_, _ = builder.WriteString(p.buildCacheConfig(ctx.Caches, r.CacheID))

	return builder.String()
}

func (p *hostConfigurationProvider) buildCacheConfig(
	caches []cache.Cache,
	cacheID *uuid.UUID,
) string {
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
	_, _ = builder.WriteString("\n")

	cacheIDNoDashes := strings.ReplaceAll(c.ID.String(), "-", "")
	nginxFprintf(&builder, "proxy_cache %s;", "cache_"+cacheIDNoDashes)

	p.appendCacheDurations(&builder, c)
	p.appendCacheMethods(&builder, c)
	p.appendCacheStandardOptions(&builder, c)
	p.appendCacheLock(&builder, c)
	p.appendCacheBypassRules(&builder, c)
	p.appendCacheFileExtensions(&builder, c)

	return builder.String()
}

func (p *hostConfigurationProvider) appendCacheDurations(
	builder *strings.Builder,
	c *cache.Cache,
) {
	for _, d := range c.Durations {
		arguments := make([]any, 0, len(d.StatusCodes)+1)
		for _, statusCode := range d.StatusCodes {
			arguments = append(arguments, statusCode)
		}
		arguments = append(arguments, directiveFragment(nginxSprintf("%ds", d.ValidTimeSeconds)))

		nginxFprintfArgs(builder, "\nproxy_cache_valid %s;", arguments...)
	}
}

func (p *hostConfigurationProvider) appendCacheMethods(
	builder *strings.Builder,
	c *cache.Cache,
) {
	if len(c.AllowedMethods) > 0 {
		methods := make([]string, len(c.AllowedMethods))
		for index, method := range c.AllowedMethods {
			methods[index] = strings.ToLower(string(method))
		}

		nginxFprintfArgs(builder, "\nproxy_cache_methods %s;", methods...)
	}
}

func (p *hostConfigurationProvider) appendCacheStandardOptions(
	builder *strings.Builder,
	c *cache.Cache,
) {
	nginxFprintf(builder, "\nproxy_cache_min_uses %d;", c.MinimumUsesBeforeCaching)
	nginxFprintf(
		builder,
		"\nproxy_cache_background_update %s;",
		directiveFragment(statusFlag(c.BackgroundUpdate)),
	)
	nginxFprintf(
		builder,
		"\nproxy_cache_revalidate %s;",
		directiveFragment(statusFlag(c.Revalidate)),
	)

	if c.IgnoreUpstreamCacheHeaders {
		_, _ = builder.WriteString("\nproxy_ignore_headers Cache-Control Expires;")
	}

	if c.CacheStatusResponseHeaderEnabled {
		_, _ = builder.WriteString("\nadd_header X-Cache-Status $upstream_cache_status;")
	}

	staleConfig := offFlag

	if len(c.UseStale) > 0 {
		staleOptions := make([]string, len(c.UseStale))
		for index, option := range c.UseStale {
			staleOptions[index] = strings.ToLower(string(option))
		}

		nginxFprintfArgs(builder, "\nproxy_cache_use_stale %s;", staleOptions...)
	} else {
		nginxFprintf(builder, "\nproxy_cache_use_stale %s;", directiveFragment(staleConfig))
	}
}

func (p *hostConfigurationProvider) appendCacheLock(builder *strings.Builder, c *cache.Cache) {
	if c.ConcurrencyLock.Enabled {
		_, _ = builder.WriteString("\nproxy_cache_lock on;")
		if c.ConcurrencyLock.TimeoutSeconds != nil {
			nginxFprintf(
				builder,
				"\nproxy_cache_lock_timeout %ds;",
				*c.ConcurrencyLock.TimeoutSeconds,
			)
		}
		if c.ConcurrencyLock.AgeSeconds != nil {
			nginxFprintf(
				builder,
				"\nproxy_cache_lock_age %ds;",
				*c.ConcurrencyLock.AgeSeconds,
			)
		}
	}
}

func (p *hostConfigurationProvider) appendCacheBypassRules(
	builder *strings.Builder,
	c *cache.Cache,
) {
	for _, rule := range c.BypassRules {
		nginxFprintf(builder, "\nproxy_cache_bypass %s;", rule)
	}

	for _, rule := range c.NoCacheRules {
		nginxFprintf(builder, "\nproxy_no_cache %s;", rule)
	}
}

func (p *hostConfigurationProvider) appendCacheFileExtensions(
	builder *strings.Builder,
	c *cache.Cache,
) {
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

	_, _ = builder.WriteString("\n")
	nginxFprintf(
		builder,
		`
			if ($uri !~* %s) {
				set $__no_cache_allowed_extension 1;
			}
			proxy_no_cache $__no_cache_allowed_extension;
		`,
		"\\.("+strings.Join(extensions, "|")+")$",
	)
}
