package host

import (
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/common/ptr"
	"dillmann.com.br/nginx-ignition/core/host"
	"dillmann.com.br/nginx-ignition/core/settings"
)

func toDto(input *host.Host, globalSettings *settings.Settings) *hostResponseDto {
	if input == nil {
		return nil
	}

	var globalBindings *[]*bindingDto
	if input.UseGlobalBindings && globalSettings != nil && len(globalSettings.GlobalBindings) > 0 {
		globalBindings = ptr.Of(toBindingDtoSlice(globalSettings.GlobalBindings))
	}

	return &hostResponseDto{
		ID:                &input.ID,
		Enabled:           &input.Enabled,
		DefaultServer:     &input.DefaultServer,
		UseGlobalBindings: &input.UseGlobalBindings,
		DomainNames:       input.DomainNames,
		Routes:            toRouteDtoSlice(input.Routes),
		Bindings:          toBindingDtoSlice(input.Bindings),
		GlobalBindings:    globalBindings,
		VPNs:              toVpnDtoSlice(input.VPNs),
		FeatureSet:        toFeatureSetDto(&input.FeatureSet),
		AccessListID:      input.AccessListID,
		CacheID:           input.CacheID,
	}
}

func toDomain(input *hostRequestDto) *host.Host {
	if input == nil {
		return nil
	}

	return &host.Host{
		Enabled:           getBoolValue(input.Enabled),
		DefaultServer:     getBoolValue(input.DefaultServer),
		UseGlobalBindings: getBoolValue(input.UseGlobalBindings),
		DomainNames:       input.DomainNames,
		Routes:            toRouteSlice(input.Routes),
		Bindings:          toBindingSlice(input.Bindings),
		VPNs:              toVpnsSlice(input.VPNs),
		FeatureSet:        *toFeatureSet(input.FeatureSet),
		AccessListID:      input.AccessListID,
		CacheID:           input.CacheID,
	}
}

func toRouteDtoSlice(routes []*host.Route) []*routeDto {
	if routes == nil {
		return nil
	}

	result := make([]*routeDto, len(routes))
	for index, route := range routes {
		result[index] = toRouteDto(route)
	}

	return result
}

func toRouteDto(route *host.Route) *routeDto {
	if route == nil {
		return nil
	}

	return &routeDto{
		Priority:     &route.Priority,
		Enabled:      &route.Enabled,
		Type:         &route.Type,
		SourcePath:   &route.SourcePath,
		Settings:     toRouteSettingsDto(&route.Settings),
		TargetUri:    route.TargetURI,
		RedirectCode: route.RedirectCode,
		Response:     toStaticResponseDto(route.Response),
		Integration:  toIntegrationConfigDto(route.Integration),
		AccessListID: route.AccessListID,
		CacheID:      route.CacheID,
		SourceCode:   toRouteSourceCodeDto(route.SourceCode),
	}
}

func toBindingDtoSlice(bindings []*host.Binding) []*bindingDto {
	if bindings == nil {
		return nil
	}

	result := make([]*bindingDto, len(bindings))
	for index, binding := range bindings {
		result[index] = toBindingDto(binding)
	}

	return result
}

func toBindingDto(binding *host.Binding) *bindingDto {
	if binding == nil {
		return nil
	}

	return &bindingDto{
		Type:          &binding.Type,
		Ip:            &binding.IP,
		Port:          &binding.Port,
		CertificateId: binding.CertificateID,
	}
}

func toVpnDtoSlice(vpns []*host.VPN) []*vpnDto {
	if vpns == nil {
		return nil
	}

	result := make([]*vpnDto, len(vpns))
	for index, vpn := range vpns {
		result[index] = toVpnDto(vpn)
	}

	return result
}

func toVpnDto(vpn *host.VPN) *vpnDto {
	if vpn == nil {
		return nil
	}

	return &vpnDto{
		VPNID: &vpn.VPNID,
		Name:  &vpn.Name,
		Host:  vpn.Host,
	}
}

func toFeatureSetDto(featureSet *host.FeatureSet) *featureSetDto {
	if featureSet == nil {
		return nil
	}

	return &featureSetDto{
		WebsocketsSupport:   &featureSet.WebsocketSupport,
		Http2Support:        &featureSet.HTTP2Support,
		RedirectHttpToHttps: &featureSet.RedirectHTTPToHTTPS,
	}
}

func toRouteSettingsDto(settings *host.RouteSettings) *routeSettingsDto {
	if settings == nil {
		return nil
	}

	return &routeSettingsDto{
		IncludeForwardHeaders:   &settings.IncludeForwardHeaders,
		ProxySslServerName:      &settings.ProxySSLServerName,
		KeepOriginalDomainName:  &settings.KeepOriginalDomainName,
		DirectoryListingEnabled: &settings.DirectoryListingEnabled,
		Custom:                  settings.Custom,
	}
}

func toStaticResponseDto(response *host.RouteStaticResponse) *staticResponseDto {
	if response == nil {
		return nil
	}

	return &staticResponseDto{
		StatusCode: &response.StatusCode,
		Payload:    response.Payload,
		Headers:    &response.Headers,
	}
}

func toIntegrationConfigDto(config *host.RouteIntegrationConfig) *integrationConfigDto {
	if config == nil {
		return nil
	}

	return &integrationConfigDto{
		IntegrationId: &config.IntegrationID,
		OptionId:      &config.OptionID,
	}
}

func toRouteSourceCodeDto(sourceCode *host.RouteSourceCode) *routeSourceCodeDto {
	if sourceCode == nil {
		return nil
	}

	return &routeSourceCodeDto{
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

func toRouteSlice(routes []*routeDto) []*host.Route {
	if routes == nil {
		return nil
	}

	result := make([]*host.Route, len(routes))
	for index, route := range routes {
		result[index] = toDomainModelRoute(route)
	}

	return result
}

func toBindingSlice(bindings []*bindingDto) []*host.Binding {
	if bindings == nil {
		return nil
	}

	result := make([]*host.Binding, len(bindings))
	for index, binding := range bindings {
		result[index] = toDomainModelBinding(binding)
	}

	return result
}

func toVpnsSlice(vpns []*vpnDto) []*host.VPN {
	if vpns == nil {
		return nil
	}

	result := make([]*host.VPN, len(vpns))
	for index, vpn := range vpns {
		result[index] = toDomainModelVPN(vpn)
	}

	return result
}

func toDomainModelRoute(input *routeDto) *host.Route {
	if input == nil {
		return nil
	}

	return &host.Route{
		Priority:     getIntValue(input.Priority),
		Enabled:      getBoolValue(input.Enabled),
		Type:         *input.Type,
		SourcePath:   getStringValue(input.SourcePath),
		Settings:     *toRouteSettings(input.Settings),
		TargetURI:    input.TargetUri,
		RedirectCode: input.RedirectCode,
		Response:     toRouteStaticResponse(input.Response),
		Integration:  toRouteIntegrationConfig(input.Integration),
		AccessListID: input.AccessListID,
		CacheID:      input.CacheID,
		SourceCode:   toRouteSourceCode(input.SourceCode),
	}
}

func toRouteSettings(input *routeSettingsDto) *host.RouteSettings {
	if input == nil {
		return nil
	}

	return &host.RouteSettings{
		IncludeForwardHeaders:   getBoolValue(input.IncludeForwardHeaders),
		ProxySSLServerName:      getBoolValue(input.ProxySslServerName),
		KeepOriginalDomainName:  getBoolValue(input.KeepOriginalDomainName),
		DirectoryListingEnabled: getBoolValue(input.DirectoryListingEnabled),
		Custom:                  input.Custom,
	}
}

func toRouteStaticResponse(input *staticResponseDto) *host.RouteStaticResponse {
	if input == nil {
		return nil
	}

	return &host.RouteStaticResponse{
		StatusCode: getIntValue(input.StatusCode),
		Payload:    input.Payload,
		Headers:    getMapValue(input.Headers),
	}
}

func toRouteIntegrationConfig(input *integrationConfigDto) *host.RouteIntegrationConfig {
	if input == nil {
		return nil
	}

	return &host.RouteIntegrationConfig{
		IntegrationID: getUuidValue(input.IntegrationId),
		OptionID:      getStringValue(input.OptionId),
	}
}

func toRouteSourceCode(input *routeSourceCodeDto) *host.RouteSourceCode {
	if input == nil {
		return nil
	}

	return &host.RouteSourceCode{
		Language:     *input.Language,
		Contents:     getStringValue(input.Code),
		MainFunction: input.MainFunction,
	}
}

func toFeatureSet(input *featureSetDto) *host.FeatureSet {
	if input == nil {
		return nil
	}

	return &host.FeatureSet{
		WebsocketSupport:    getBoolValue(input.WebsocketsSupport),
		HTTP2Support:        getBoolValue(input.Http2Support),
		RedirectHTTPToHTTPS: getBoolValue(input.RedirectHttpToHttps),
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

func getUuidValue(value *uuid.UUID) uuid.UUID {
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

func toDomainModelBinding(input *bindingDto) *host.Binding {
	if input == nil {
		return nil
	}

	return &host.Binding{
		Type:          *input.Type,
		IP:            getStringValue(input.Ip),
		Port:          getIntValue(input.Port),
		CertificateID: input.CertificateId,
	}
}

func toDomainModelVPN(input *vpnDto) *host.VPN {
	if input == nil {
		return nil
	}

	return &host.VPN{
		VPNID: getUuidValue(input.VPNID),
		Name:  getStringValue(input.Name),
		Host:  input.Host,
	}
}
