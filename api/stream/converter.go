package stream

import (
	"dillmann.com.br/nginx-ignition/core/common/ptr"
	"dillmann.com.br/nginx-ignition/core/stream"
)

func toDTO(input *stream.Stream) *streamResponseDTO {
	if input == nil {
		return nil
	}

	return &streamResponseDTO{
		ID:             &input.ID,
		Enabled:        &input.Enabled,
		Name:           &input.Name,
		Type:           ptr.Of(string(input.Type)),
		FeatureSet:     toFeatureSetDTO(&input.FeatureSet),
		DefaultBackend: toBackendDTO(&input.DefaultBackend),
		Binding:        toAddressDTO(&input.Binding),
		Routes:         toRouteDTOs(input.Routes),
	}
}

func toDomain(input *streamRequestDTO) *stream.Stream {
	if input == nil {
		return nil
	}

	return &stream.Stream{
		Enabled:        getBoolValue(input.Enabled),
		Name:           getStringValue(input.Name),
		Type:           stream.Type(getStringValue(input.Type)),
		FeatureSet:     *toFeatureSet(input.FeatureSet),
		DefaultBackend: *toBackend(input.DefaultBackend),
		Binding:        *toAddress(input.Binding),
		Routes:         toRoutes(input.Routes),
	}
}

func toFeatureSetDTO(featureSet *stream.FeatureSet) *featureSetDTO {
	if featureSet == nil {
		return nil
	}

	return &featureSetDTO{
		UseProxyProtocol: &featureSet.UseProxyProtocol,
		SocketKeepAlive:  &featureSet.SocketKeepAlive,
		TCPKeepAlive:     &featureSet.TCPKeepAlive,
		TCPNoDelay:       &featureSet.TCPNoDelay,
		TCPDeferred:      &featureSet.TCPDeferred,
	}
}

func toAddressDTO(address *stream.Address) *addressDTO {
	if address == nil {
		return nil
	}

	return &addressDTO{
		Protocol: address.Protocol,
		Address:  &address.Address,
		Port:     address.Port,
	}
}

func toBackend(backend *backendDTO) *stream.Backend {
	if backend == nil {
		return nil
	}

	return &stream.Backend{
		Weight:         backend.Weight,
		Address:        *toAddress(backend.Target),
		CircuitBreaker: toCircuitBreaker(backend.CircuitBreaker),
	}
}

func toCircuitBreaker(input *circuitBreakerDTO) *stream.CircuitBreaker {
	if input == nil {
		return nil
	}

	return &stream.CircuitBreaker{
		MaxFailures: getIntValue(input.MaxFailures),
		OpenSeconds: getIntValue(input.OpenSeconds),
	}
}

func toCircuitBreakerDTO(input *stream.CircuitBreaker) *circuitBreakerDTO {
	if input == nil {
		return nil
	}

	return &circuitBreakerDTO{
		MaxFailures: &input.MaxFailures,
		OpenSeconds: &input.OpenSeconds,
	}
}

func toBackendDTO(input *stream.Backend) *backendDTO {
	if input == nil {
		return nil
	}

	return &backendDTO{
		Weight:         input.Weight,
		Target:         toAddressDTO(&input.Address),
		CircuitBreaker: toCircuitBreakerDTO(input.CircuitBreaker),
	}
}

func toRouteDTOs(input []stream.Route) []routeDTO {
	output := make([]routeDTO, len(input))
	for index := range input {
		output[index] = routeDTO{
			DomainNames: input[index].DomainNames,
			Backends:    toBackendDTOs(input[index].Backends),
		}
	}

	return output
}

func toBackendDTOs(input []stream.Backend) []backendDTO {
	output := make([]backendDTO, len(input))
	for index := range input {
		output[index] = *toBackendDTO(&input[index])
	}

	return output
}

func toRoutes(input []routeDTO) []stream.Route {
	output := make([]stream.Route, len(input))
	for index := range input {
		output[index] = stream.Route{
			DomainNames: input[index].DomainNames,
			Backends:    toBackends(input[index].Backends),
		}
	}

	return output
}

func toBackends(input []backendDTO) []stream.Backend {
	output := make([]stream.Backend, len(input))
	for index := range input {
		output[index] = *toBackend(&input[index])
	}

	return output
}

func toFeatureSet(input *featureSetDTO) *stream.FeatureSet {
	if input == nil {
		return nil
	}

	return &stream.FeatureSet{
		UseProxyProtocol: getBoolValue(input.UseProxyProtocol),
		SocketKeepAlive:  getBoolValue(input.SocketKeepAlive),
		TCPKeepAlive:     getBoolValue(input.TCPKeepAlive),
		TCPNoDelay:       getBoolValue(input.TCPNoDelay),
		TCPDeferred:      getBoolValue(input.TCPDeferred),
	}
}

func toAddress(input *addressDTO) *stream.Address {
	if input == nil {
		return nil
	}

	return &stream.Address{
		Protocol: input.Protocol,
		Address:  getStringValue(input.Address),
		Port:     input.Port,
	}
}

func getBoolValue(value *bool) bool {
	if value == nil {
		return false
	}

	return *value
}

func getStringValue(value *string) string {
	if value == nil {
		return ""
	}

	return *value
}

func getIntValue(value *int) int {
	if value == nil {
		return 0
	}

	return *value
}
