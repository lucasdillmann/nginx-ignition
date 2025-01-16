package host

import (
	"dillmann.com.br/nginx-ignition/core/common/validation"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

type validator struct {
	hostRepository *Repository
	delegate       *validation.ConsistencyValidator
}

func newValidator(hostRepository *Repository) *validator {
	return &validator{
		hostRepository: hostRepository,
		delegate:       validation.NewValidator(),
	}
}

const (
	bindingsPath              = "bindings"
	minimumPort               = 1
	maximumPort               = 65535
	minimumRedirectStatusCode = 300
	maximumRedirectStatusCode = 399
	minimumStatusCode         = 100
	maximumStatusCode         = 599
)

var tldPattern = regexp.MustCompile(`(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\\.)+[a-z0-9][a-z0-9-]{0,61}[a-z0-9]`)

func (v *validator) validate(host *Host) error {
	if err := v.validateDefaultFlag(host); err != nil {
		return err
	}

	if err := v.validateBindings(host); err != nil {
		return err
	}

	v.validateDomainNames(host)
	v.validateRoutes(host)

	return v.delegate.Result()
}

func (v *validator) validateDefaultFlag(host *Host) error {
	if !host.DefaultServer {
		return nil
	}

	defaultServer, err := (*v.hostRepository).FindDefault()
	if err != nil {
		return err
	}

	if defaultServer != nil && host.ID != defaultServer.ID {
		v.delegate.Add("defaultServer", "There's already another host marked as the default one")
	}

	if len(host.DomainNames) > 0 {
		v.delegate.Add("domainNames", "Must be empty when the host is the default one")
	}

	return nil
}

func (v *validator) validateDomainNames(host *Host) {
	if len(host.DomainNames) == 0 && !host.DefaultServer {
		v.delegate.Add("domainNames", "At least one domain name must be informed")
	}

	for index, domainName := range host.DomainNames {
		if !tldPattern.MatchString(*domainName) {
			v.delegate.Add("domainNames["+strconv.Itoa(index)+"]", "Value is not a valid domain name")
		}
	}
}

func (v *validator) validateBindings(host *Host) error {
	if host.UseGlobalBindings && len(host.Bindings) > 0 {
		v.delegate.Add(bindingsPath, "Must be empty when using global bindings")
	}

	if !host.UseGlobalBindings {
		if len(host.Bindings) == 0 {
			v.delegate.Add(bindingsPath, "At least one binding must be informed")
		}

		for index, binding := range host.Bindings {
			if err := v.validateBinding(bindingsPath, binding, index); err != nil {
				return err
			}
		}
	}

	return nil
}

func (v *validator) validateBinding(pathPrefix string, binding *Binding, index int) error {
	if net.ParseIP(binding.IP) == nil {
		v.delegate.Add(pathPrefix+"["+strconv.Itoa(index)+"].ip", "Not a valid IPv4 or IPv6 address")
	}

	if binding.Port < minimumPort || binding.Port > maximumPort {
		v.delegate.Add(
			pathPrefix+"["+strconv.Itoa(index)+"].port",
			"Value must be between "+strconv.Itoa(minimumPort)+" and "+strconv.Itoa(maximumPort),
		)
	}

	certificateIdField := pathPrefix + "[" + strconv.Itoa(index) + "].certificateId"

	switch {
	case binding.Type == HttpBindingType && binding.CertificateID != nil:
		v.delegate.Add(certificateIdField, "Value cannot be informed for a HTTP binding")
	case binding.Type == HttpsBindingType && binding.CertificateID == nil:
		v.delegate.Add(certificateIdField, "Value must be informed for a HTTPS binding")
	case binding.Type == HttpsBindingType:
		exists, err := (*v.hostRepository).ExistsCertificateByID(*binding.CertificateID)
		if err != nil {
			return err
		}

		if !exists {
			v.delegate.Add(certificateIdField, "No SSL certificate found with provided ID")
		}
	default:
		v.delegate.Add(pathPrefix+"["+strconv.Itoa(index)+"].type", "Invalid value")
	}

	return nil
}

func (v *validator) validateRoutes(host *Host) {
	if len(host.Routes) == 0 {
		v.delegate.Add("routes", "At least one route must be informed")
	}

	priorityMap := make(map[int]int)
	for _, route := range host.Routes {
		priorityMap[route.Priority]++
	}

	for priority, count := range priorityMap {
		if count > 1 {
			v.delegate.Add(
				"routes",
				"Priority "+strconv.Itoa(priority)+" is duplicated in two or more routes",
			)
		}
	}

	distinctPaths := make(map[string]bool)
	for index, route := range host.Routes {
		v.validateRoute(route, index, &distinctPaths)
	}
}

func (v *validator) validateRoute(route *Route, index int, distinctPaths *map[string]bool) {
	if (*distinctPaths)[route.SourcePath] {
		v.delegate.Add(
			"routes["+strconv.Itoa(index)+"].sourcePath",
			"Source path was already used in another route",
		)
	} else {
		(*distinctPaths)[route.SourcePath] = true
	}

	switch route.Type {
	case ProxyRouteType:
		v.validateProxyRoute(route, index)
	case RedirectRouteType:
		v.validateRedirectRoute(route, index)
	case StaticResponseRouteType:
		v.validateStaticResponseRoute(route, index)
	case IntegrationRouteType:
		v.validateIntegrationRoute(route, index)
	case SourceCodeRouteType:
		v.validateSourceCodeRoute(route, index)
	default:
		v.delegate.Add("routes["+strconv.Itoa(index)+"].type", "Invalid value")
	}
}

func (v *validator) validateProxyRoute(route *Route, index int) {
	targetUriField := "routes[" + strconv.Itoa(index) + "].targetUri"
	if route.TargetURI == nil || strings.TrimSpace(*route.TargetURI) == "" {
		v.delegate.Add(targetUriField, "Value is required when the type of the route is proxy")
	} else {
		if _, err := url.ParseRequestURI(*route.TargetURI); err != nil {
			v.delegate.Add(targetUriField, "Value is not a valid URL")
		}
	}
}

func (v *validator) validateRedirectRoute(route *Route, index int) {
	targetUriField := "routes[" + strconv.Itoa(index) + "].targetUri"
	if route.TargetURI == nil || strings.TrimSpace(*route.TargetURI) == "" {
		v.delegate.Add(targetUriField, "Value is required when the type of the route is redirect")
	} else {
		if _, err := url.ParseRequestURI(*route.TargetURI); err != nil {
			v.delegate.Add(targetUriField, "Value is not a valid URI")
		}
	}

	if route.RedirectCode == nil ||
		*route.RedirectCode < minimumRedirectStatusCode ||
		*route.RedirectCode > maximumRedirectStatusCode {
		v.delegate.Add(
			"routes["+strconv.Itoa(index)+"].redirectCode",
			"Value must be between "+strconv.Itoa(minimumRedirectStatusCode)+" and "+strconv.Itoa(maximumRedirectStatusCode),
		)
	}
}

func (v *validator) validateStaticResponseRoute(route *Route, index int) {
	if route.Response == nil {
		v.delegate.Add(
			"routes["+strconv.Itoa(index)+"].response",
			"A value is required when the type of the route is static response",
		)
		return
	}

	if route.Response.StatusCode < minimumStatusCode || route.Response.StatusCode > maximumStatusCode {
		v.delegate.Add(
			"routes["+strconv.Itoa(index)+"].response.statusCode",
			"Value must be between "+strconv.Itoa(minimumStatusCode)+" and "+strconv.Itoa(maximumStatusCode),
		)
	}
}

func (v *validator) validateIntegrationRoute(route *Route, index int) {
	requiredMessage := "Value is required when the type of the route is integration"

	if route.Integration == nil {
		v.delegate.Add("routes["+strconv.Itoa(index)+"].integration", requiredMessage)
		return
	}

	if strings.TrimSpace(route.Integration.IntegrationID) == "" {
		v.delegate.Add("routes["+strconv.Itoa(index)+"].integration.integrationId", requiredMessage)
	}

	if strings.TrimSpace(route.Integration.OptionID) == "" {
		v.delegate.Add("routes["+strconv.Itoa(index)+"].integration.optionId", requiredMessage)
	}
}

func (v *validator) validateSourceCodeRoute(route *Route, index int) {
	requiredMessage := "Value is required when the type of the route is source code"

	if route.SourceCode == nil {
		v.delegate.Add("routes["+strconv.Itoa(index)+"].sourceCode", requiredMessage)
		return
	}

	if strings.TrimSpace(route.SourceCode.Contents) == "" {
		v.delegate.Add("routes["+strconv.Itoa(index)+"].sourceCode.code", requiredMessage)
	}

	if route.SourceCode.Language != JavascriptCodeLanguage && route.SourceCode.Language != LuaCodeLanguage {
		v.delegate.Add(
			"routes["+strconv.Itoa(index)+"].sourceCode.language",
			"Invalid value",
		)
	}

	if route.SourceCode.Language == JavascriptCodeLanguage &&
		(route.SourceCode.MainFunction == nil || strings.TrimSpace(*route.SourceCode.MainFunction) == "") {
		v.delegate.Add(
			"routes["+strconv.Itoa(index)+"].sourceCode.mainFunction",
			"Value is required when the language is JavaScript",
		)
	}
}
