package host

import (
	"dillmann.com.br/nginx-ignition/core/host"
	"encoding/json"
)

func toDomain(model *hostModel) (*host.Host, error) {
	bindings := make([]host.Binding, len(model.Bindings))
	for index, binding := range model.Bindings {
		bindings[index] = host.Binding{
			ID:            binding.ID,
			Type:          host.BindingType(binding.Type),
			IP:            binding.IP,
			Port:          binding.Port,
			CertificateID: binding.CertificateID,
		}
	}

	routes := make([]host.Route, len(model.Routes))
	for index, route := range model.Routes {
		headers, err := parseHeaders(route.StaticResponseHeaders)
		if err != nil {
			return nil, err
		}

		var response *host.RouteStaticResponse
		if route.StaticResponseCode != nil {
			response = &host.RouteStaticResponse{
				StatusCode: *route.StaticResponseCode,
				Headers:    headers,
				Payload:    route.StaticResponsePayload,
			}
		}

		var integration *host.RouteIntegrationConfig
		if route.IntegrationOptionID != nil {
			integration = &host.RouteIntegrationConfig{
				IntegrationID: *route.IntegrationID,
				OptionID:      *route.IntegrationOptionID,
			}
		}

		var sourceCode *host.RouteSourceCode
		if route.CodeLanguage != nil {
			sourceCode = &host.RouteSourceCode{
				Language:     host.CodeLanguage(*route.CodeLanguage),
				Contents:     *route.CodeContents,
				MainFunction: route.CodeMainFunction,
			}
		}

		routes[index] = host.Route{
			ID:           route.ID,
			Priority:     route.Priority,
			Enabled:      route.Enabled,
			Type:         host.RouteType(route.Type),
			SourcePath:   route.SourcePath,
			TargetURI:    route.TargetURI,
			RedirectCode: route.RedirectCode,
			AccessListID: route.AccessListID,
			Settings: host.RouteSettings{
				IncludeForwardHeaders:  route.IncludeForwardHeaders,
				ProxySSLServerName:     route.ProxySSLServerName,
				KeepOriginalDomainName: route.KeepOriginalDomainName,
				Custom:                 route.CustomSettings,
			},
			Response:    response,
			Integration: integration,
			SourceCode:  sourceCode,
		}
	}

	return &host.Host{
		ID:                model.ID,
		Enabled:           model.Enabled,
		DefaultServer:     model.DefaultServer,
		UseGlobalBindings: model.UseGlobalBindings,
		DomainNames:       model.DomainNames,
		Routes:            routes,
		Bindings:          bindings,
		FeatureSet: host.FeatureSet{
			WebsocketSupport:    model.WebsocketSupport,
			HTTP2Support:        model.HTTP2Support,
			RedirectHTTPToHTTPS: model.RedirectHTTPToHTTPS,
		},
		AccessListID: model.AccessListID,
	}, nil
}

func toModel(domain *host.Host) (*hostModel, error) {
	bindings := make([]hostBindingModel, len(domain.Bindings))
	for index, binding := range domain.Bindings {
		bindings[index] = hostBindingModel{
			ID:            binding.ID,
			HostID:        domain.ID,
			Type:          string(binding.Type),
			IP:            binding.IP,
			Port:          binding.Port,
			CertificateID: binding.CertificateID,
		}
	}

	routes := make([]hostRouteModel, len(domain.Routes))
	for index, route := range domain.Routes {
		headers, err := formatHeaders(route.Response.Headers)
		if err != nil {
			return nil, err
		}

		var integrationID, integrationOptionID *string
		if route.Integration != nil {
			integrationID = &route.Integration.IntegrationID
			integrationOptionID = &route.Integration.OptionID
		}

		var codeLanguage, codeContents *string
		if route.SourceCode != nil {
			codeLanguage = (*string)(&route.SourceCode.Language)
			codeContents = &route.SourceCode.Contents
		}

		routes[index] = hostRouteModel{
			ID:                     route.ID,
			HostID:                 domain.ID,
			Priority:               route.Priority,
			Type:                   string(route.Type),
			SourcePath:             route.SourcePath,
			TargetURI:              route.TargetURI,
			CustomSettings:         route.Settings.Custom,
			StaticResponseCode:     &route.Response.StatusCode,
			StaticResponsePayload:  route.Response.Payload,
			StaticResponseHeaders:  headers,
			RedirectCode:           route.RedirectCode,
			IntegrationID:          integrationID,
			IntegrationOptionID:    integrationOptionID,
			IncludeForwardHeaders:  route.Settings.IncludeForwardHeaders,
			ProxySSLServerName:     route.Settings.ProxySSLServerName,
			KeepOriginalDomainName: route.Settings.KeepOriginalDomainName,
			AccessListID:           route.AccessListID,
			CodeLanguage:           codeLanguage,
			CodeContents:           codeContents,
			CodeMainFunction:       route.SourceCode.MainFunction,
			Enabled:                route.Enabled,
		}
	}

	return &hostModel{
		ID:                  domain.ID,
		Enabled:             domain.Enabled,
		DefaultServer:       domain.DefaultServer,
		DomainNames:         domain.DomainNames,
		WebsocketSupport:    domain.FeatureSet.WebsocketSupport,
		HTTP2Support:        domain.FeatureSet.HTTP2Support,
		RedirectHTTPToHTTPS: domain.FeatureSet.RedirectHTTPToHTTPS,
		UseGlobalBindings:   domain.UseGlobalBindings,
		AccessListID:        domain.AccessListID,
		Bindings:            bindings,
		Routes:              routes,
	}, nil
}

func parseHeaders(headers *string) (map[string]string, error) {
	if headers == nil {
		return nil, nil
	}
	var result map[string]string
	if err := json.Unmarshal([]byte(*headers), &result); err != nil {
		return nil, err
	}
	return result, nil
}

func formatHeaders(headers map[string]string) (*string, error) {
	if headers == nil {
		return nil, nil
	}
	result, err := json.Marshal(headers)
	if err != nil {
		return nil, err
	}
	strResult := string(result)
	return &strResult, nil
}
