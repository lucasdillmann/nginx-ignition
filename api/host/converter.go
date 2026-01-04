package host

import (
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/binding"
	"dillmann.com.br/nginx-ignition/core/host"
	"dillmann.com.br/nginx-ignition/core/settings"
)

func toDTO(input *host.Host, globalSettings *settings.Settings) *hostResponseDTO {
	if input == nil {
		return nil
	}

	globalBindings := make([]bindingDTO, 0)
	if input.UseGlobalBindings && globalSettings != nil && len(globalSettings.GlobalBindings) > 0 {
		globalBindings = toBindingDTOSlice(globalSettings.GlobalBindings)
	}

	return &hostResponseDTO{
		ID:                &input.ID,
		Enabled:           &input.Enabled,
		DefaultServer:     &input.DefaultServer,
		UseGlobalBindings: &input.UseGlobalBindings,
		DomainNames:       input.DomainNames,
		Routes:            toRouteDTOSlice(input.Routes),
		Bindings:          toBindingDTOSlice(input.Bindings),
		GlobalBindings:    globalBindings,
		VPNs:              toVpnDTOSlice(input.VPNs),
		FeatureSet:        toFeatureSetDTO(&input.FeatureSet),
		AccessListID:      input.AccessListID,
		CacheID:           input.CacheID,
	}
}

func toDomain(input *hostRequestDTO) *host.Host {
	if input == nil {
		return nil
	}

	var featureSet host.FeatureSet
	if fs := toFeatureSet(input.FeatureSet); fs != nil {
		featureSet = *fs
	}

	return &host.Host{
		Enabled:           getBoolValue(input.Enabled),
		DefaultServer:     getBoolValue(input.DefaultServer),
		UseGlobalBindings: getBoolValue(input.UseGlobalBindings),
		DomainNames:       input.DomainNames,
		Routes:            toRouteSlice(input.Routes),
		Bindings:          toBindingSlice(input.Bindings),
		VPNs:              toVPNsSlice(input.VPNs),
		FeatureSet:        featureSet,
		AccessListID:      input.AccessListID,
		CacheID:           input.CacheID,
	}
}

func toRouteDTOSlice(routes []host.Route) []routeDTO {
	result := make([]routeDTO, len(routes))
	for index, route := range routes {
		result[index] = toRouteDTO(&route)
	}

	return result
}

func toRouteDTO(route *host.Route) routeDTO {
	return routeDTO{
		Priority:     &route.Priority,
		Enabled:      &route.Enabled,
		Type:         &route.Type,
		SourcePath:   &route.SourcePath,
		Settings:     toRouteSettingsDTO(&route.Settings),
		TargetURI:    route.TargetURI,
		RedirectCode: route.RedirectCode,
		Response:     toStaticResponseDTO(route.Response),
		Integration:  toIntegrationConfigDTO(route.Integration),
		AccessListID: route.AccessListID,
		CacheID:      route.CacheID,
		SourceCode:   toRouteSourceCodeDTO(route.SourceCode),
	}
}

func toBindingDTOSlice(bindings []binding.Binding) []bindingDTO {
	result := make([]bindingDTO, len(bindings))
	for index, b := range bindings {
		result[index] = bindingDTO{
			Type:          &b.Type,
			IP:            &b.IP,
			Port:          &b.Port,
			CertificateID: b.CertificateID,
		}
	}

	return result
}

func toVpnDTOSlice(vpns []host.VPN) []vpnDTO {
	result := make([]vpnDTO, len(vpns))
	for index, vpn := range vpns {
		result[index] = vpnDTO{
			VPNID: &vpn.VPNID,
			Name:  &vpn.Name,
			Host:  vpn.Host,
		}
	}

	return result
}

func toFeatureSetDTO(featureSet *host.FeatureSet) *featureSetDTO {
	if featureSet == nil {
		return nil
	}

	return &featureSetDTO{
		WebsocketsSupport:   &featureSet.WebsocketSupport,
		HTTP2Support:        &featureSet.HTTP2Support,
		RedirectHTTPToHTTPS: &featureSet.RedirectHTTPToHTTPS,
	}
}

func toRouteSettingsDTO(set *host.RouteSettings) *routeSettingsDTO {
	if set == nil {
		return nil
	}

	return &routeSettingsDTO{
		IncludeForwardHeaders:   &set.IncludeForwardHeaders,
		ProxySslServerName:      &set.ProxySSLServerName,
		KeepOriginalDomainName:  &set.KeepOriginalDomainName,
		DirectoryListingEnabled: &set.DirectoryListingEnabled,
		Custom:                  set.Custom,
	}
}

func toStaticResponseDTO(response *host.RouteStaticResponse) *staticResponseDTO {
	if response == nil {
		return nil
	}

	return &staticResponseDTO{
		StatusCode: &response.StatusCode,
		Payload:    response.Payload,
		Headers:    &response.Headers,
	}
}

func toIntegrationConfigDTO(config *host.RouteIntegrationConfig) *integrationConfigDTO {
	if config == nil {
		return nil
	}

	return &integrationConfigDTO{
		IntegrationID: &config.IntegrationID,
		OptionID:      &config.OptionID,
	}
}

func toRouteSourceCodeDTO(sourceCode *host.RouteSourceCode) *routeSourceCodeDTO {
	if sourceCode == nil {
		return nil
	}

	return &routeSourceCodeDTO{
		Language:     &sourceCode.Language,
		Code:         &sourceCode.Contents,
		MainFunction: sourceCode.MainFunction,
	}
}

func getBoolValue(value *bool) bool {
	if value == nil {
		return false
	}
	return *value
}

func toRouteSlice(routes []routeDTO) []host.Route {
	result := make([]host.Route, len(routes))
	for index, route := range routes {
		result[index] = host.Route{
			Priority:     getIntValue(route.Priority),
			Enabled:      getBoolValue(route.Enabled),
			Type:         *route.Type,
			SourcePath:   getStringValue(route.SourcePath),
			Settings:     toRouteSettings(route.Settings),
			TargetURI:    route.TargetURI,
			RedirectCode: route.RedirectCode,
			Response:     toRouteStaticResponse(route.Response),
			Integration:  toRouteIntegrationConfig(route.Integration),
			AccessListID: route.AccessListID,
			CacheID:      route.CacheID,
			SourceCode:   toRouteSourceCode(route.SourceCode),
		}
	}

	return result
}

func toBindingSlice(bindings []bindingDTO) []binding.Binding {
	result := make([]binding.Binding, len(bindings))
	for index, b := range bindings {
		result[index] = binding.Binding{
			Type:          *b.Type,
			IP:            getStringValue(b.IP),
			Port:          getIntValue(b.Port),
			CertificateID: b.CertificateID,
		}
	}

	return result
}

func toVPNsSlice(vpns []vpnDTO) []host.VPN {
	result := make([]host.VPN, len(vpns))
	for index, vpn := range vpns {
		result[index] = host.VPN{
			VPNID: getUUIDValue(vpn.VPNID),
			Name:  getStringValue(vpn.Name),
			Host:  vpn.Host,
		}
	}

	return result
}

func toRouteSettings(input *routeSettingsDTO) host.RouteSettings {
	return host.RouteSettings{
		IncludeForwardHeaders:   getBoolValue(input.IncludeForwardHeaders),
		ProxySSLServerName:      getBoolValue(input.ProxySslServerName),
		KeepOriginalDomainName:  getBoolValue(input.KeepOriginalDomainName),
		DirectoryListingEnabled: getBoolValue(input.DirectoryListingEnabled),
		Custom:                  input.Custom,
	}
}

func toRouteStaticResponse(input *staticResponseDTO) *host.RouteStaticResponse {
	if input == nil {
		return nil
	}

	return &host.RouteStaticResponse{
		StatusCode: getIntValue(input.StatusCode),
		Payload:    input.Payload,
		Headers:    getMapValue(input.Headers),
	}
}

func toRouteIntegrationConfig(input *integrationConfigDTO) *host.RouteIntegrationConfig {
	if input == nil {
		return nil
	}

	return &host.RouteIntegrationConfig{
		IntegrationID: getUUIDValue(input.IntegrationID),
		OptionID:      getStringValue(input.OptionID),
	}
}

func toRouteSourceCode(input *routeSourceCodeDTO) *host.RouteSourceCode {
	if input == nil {
		return nil
	}

	return &host.RouteSourceCode{
		Language:     *input.Language,
		Contents:     getStringValue(input.Code),
		MainFunction: input.MainFunction,
	}
}

func toFeatureSet(input *featureSetDTO) *host.FeatureSet {
	if input == nil {
		return nil
	}

	return &host.FeatureSet{
		WebsocketSupport:    getBoolValue(input.WebsocketsSupport),
		HTTP2Support:        getBoolValue(input.HTTP2Support),
		RedirectHTTPToHTTPS: getBoolValue(input.RedirectHTTPToHTTPS),
	}
}

func getIntValue(value *int) int {
	if value == nil {
		return 0
	}
	return *value
}

func getStringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func getUUIDValue(value *uuid.UUID) uuid.UUID {
	if value == nil {
		return uuid.Nil
	}

	return *value
}

func getMapValue(value *map[string]string) map[string]string {
	if value == nil {
		return nil
	}

	return *value
}
