package host

import (
	"encoding/json"

	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/common/pointers"
	"dillmann.com.br/nginx-ignition/core/host"
)

func toDomain(model *hostModel) (*host.Host, error) {
	bindings := make([]*host.Binding, len(model.Bindings))
	for index, binding := range model.Bindings {
		bindings[index] = &host.Binding{
			ID:            binding.ID,
			Type:          host.BindingType(binding.Type),
			IP:            binding.IP,
			Port:          binding.Port,
			CertificateID: binding.CertificateID,
		}
	}

	vpns := make([]*host.VPN, len(model.VPNs))
	for index, vpn := range model.VPNs {
		vpns[index] = &host.VPN{
			VPNID: vpn.VPNID,
			Name:  vpn.Name,
		}
	}

	routes := make([]*host.Route, len(model.Routes))
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

		routes[index] = &host.Route{
			ID:           route.ID,
			Priority:     route.Priority,
			Enabled:      route.Enabled,
			Type:         host.RouteType(route.Type),
			SourcePath:   route.SourcePath,
			TargetURI:    route.TargetURI,
			RedirectCode: route.RedirectCode,
			AccessListID: route.AccessListID,
			Settings: host.RouteSettings{
				IncludeForwardHeaders:   route.IncludeForwardHeaders,
				ProxySSLServerName:      route.ProxySSLServerName,
				KeepOriginalDomainName:  route.KeepOriginalDomainName,
				DirectoryListingEnabled: route.DirectoryListingEnabled,
				Custom:                  route.CustomSettings,
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
		DomainNames:       pointers.Reference(model.DomainNames),
		Routes:            routes,
		Bindings:          bindings,
		VPNs:              vpns,
		FeatureSet: host.FeatureSet{
			WebsocketSupport:    model.WebsocketSupport,
			HTTP2Support:        model.HTTP2Support,
			RedirectHTTPToHTTPS: model.RedirectHTTPToHTTPS,
		},
		AccessListID: model.AccessListID,
	}, nil
}

func toModel(domain *host.Host) (*hostModel, error) {
	bindings := make([]*hostBindingModel, len(domain.Bindings))
	for index, binding := range domain.Bindings {
		bindings[index] = &hostBindingModel{
			ID:            binding.ID,
			HostID:        domain.ID,
			Type:          string(binding.Type),
			IP:            binding.IP,
			Port:          binding.Port,
			CertificateID: binding.CertificateID,
		}
	}

	vpns := make([]*hostVpnModel, len(domain.VPNs))
	for index, vpn := range domain.VPNs {
		vpns[index] = &hostVpnModel{
			HostID: domain.ID,
			VPNID:  vpn.VPNID,
			Name:   vpn.Name,
		}
	}

	routes := make([]*hostRouteModel, len(domain.Routes))
	for index, route := range domain.Routes {
		var responseHeaders, responsePayload *string
		var responseStatusCode *int

		if route.Response != nil {
			var err error
			responseHeaders, err = formatHeaders(route.Response.Headers)
			if err != nil {
				return nil, err
			}

			responsePayload = route.Response.Payload
			responseStatusCode = &route.Response.StatusCode
		}

		var integrationID *uuid.UUID
		var integrationOptionID *string
		if route.Integration != nil {
			integrationID = &route.Integration.IntegrationID
			integrationOptionID = &route.Integration.OptionID
		}

		var codeLanguage, codeContents, codeMainFunction *string
		if route.SourceCode != nil {
			codeLanguage = (*string)(&route.SourceCode.Language)
			codeContents = &route.SourceCode.Contents
			codeMainFunction = route.SourceCode.MainFunction
		}

		routes[index] = &hostRouteModel{
			ID:                      route.ID,
			HostID:                  domain.ID,
			Priority:                route.Priority,
			Type:                    string(route.Type),
			SourcePath:              route.SourcePath,
			TargetURI:               route.TargetURI,
			CustomSettings:          route.Settings.Custom,
			StaticResponseCode:      responseStatusCode,
			StaticResponsePayload:   responsePayload,
			StaticResponseHeaders:   responseHeaders,
			RedirectCode:            route.RedirectCode,
			IntegrationID:           integrationID,
			IntegrationOptionID:     integrationOptionID,
			IncludeForwardHeaders:   route.Settings.IncludeForwardHeaders,
			ProxySSLServerName:      route.Settings.ProxySSLServerName,
			KeepOriginalDomainName:  route.Settings.KeepOriginalDomainName,
			DirectoryListingEnabled: route.Settings.DirectoryListingEnabled,
			AccessListID:            route.AccessListID,
			CodeLanguage:            codeLanguage,
			CodeContents:            codeContents,
			CodeMainFunction:        codeMainFunction,
			Enabled:                 route.Enabled,
		}
	}

	return &hostModel{
		ID:                  domain.ID,
		Enabled:             domain.Enabled,
		DefaultServer:       domain.DefaultServer,
		DomainNames:         pointers.Dereference(domain.DomainNames),
		WebsocketSupport:    domain.FeatureSet.WebsocketSupport,
		HTTP2Support:        domain.FeatureSet.HTTP2Support,
		RedirectHTTPToHTTPS: domain.FeatureSet.RedirectHTTPToHTTPS,
		UseGlobalBindings:   domain.UseGlobalBindings,
		AccessListID:        domain.AccessListID,
		Bindings:            bindings,
		Routes:              routes,
		VPNs:                vpns,
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
