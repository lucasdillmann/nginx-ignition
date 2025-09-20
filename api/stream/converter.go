package stream

import (
	"github.com/aws/smithy-go/ptr"

	"dillmann.com.br/nginx-ignition/core/stream"
)

func toDto(input *stream.Stream) *streamResponseDto {
	if input == nil {
		return nil
	}

	return &streamResponseDto{
		ID:             &input.ID,
		Enabled:        &input.Enabled,
		Name:           &input.Name,
		Type:           ptr.String(string(input.Type)),
		FeatureSet:     toFeatureSetDto(&input.FeatureSet),
		DefaultBackend: toBackendDto(&input.DefaultBackend),
		Binding:        toAddressDto(&input.Binding),
		Routes:         toRouteDtos(input.Routes),
	}
}

func toDomain(input *streamRequestDto) *stream.Stream {
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

func toFeatureSetDto(featureSet *stream.FeatureSet) *featureSetDto {
	if featureSet == nil {
		return nil
	}

	return &featureSetDto{
		UseProxyProtocol: &featureSet.UseProxyProtocol,
		SocketKeepAlive:  &featureSet.SocketKeepAlive,
		TCPKeepAlive:     &featureSet.TCPKeepAlive,
		TCPNoDelay:       &featureSet.TCPNoDelay,
		TCPDeferred:      &featureSet.TCPDeferred,
	}
}

func toAddressDto(address *stream.Address) *addressDto {
	if address == nil {
		return nil
	}

	return &addressDto{
		Protocol: address.Protocol,
		Address:  &address.Address,
		Port:     address.Port,
	}
}

func toBackend(backend *backendDto) *stream.Backend {
	if backend == nil {
		return nil
	}
	return &stream.Backend{
		Weight:         backend.Weight,
		Address:        *toAddress(backend.Target),
		CircuitBreaker: toCircuitBreaker(backend.CircuitBreaker),
	}
}

func toCircuitBreaker(input *circuitBreakerDto) *stream.CircuitBreaker {
	if input == nil {
		return nil
	}
	return &stream.CircuitBreaker{
		MaxFailures: getIntValue(input.MaxFailures),
		OpenSeconds: getIntValue(input.OpenSeconds),
	}
}

func toCircuitBreakerDto(input *stream.CircuitBreaker) *circuitBreakerDto {
	if input == nil {
		return nil
	}
	return &circuitBreakerDto{
		MaxFailures: &input.MaxFailures,
		OpenSeconds: &input.OpenSeconds,
	}
}

func toBackendDto(input *stream.Backend) *backendDto {
	if input == nil {
		return nil
	}
	return &backendDto{
		Weight:         input.Weight,
		Target:         toAddressDto(&input.Address),
		CircuitBreaker: toCircuitBreakerDto(input.CircuitBreaker),
	}
}

func toRouteDtos(input []stream.Route) *[]routeDto {
	if input == nil {
		return nil
	}

	output := make([]routeDto, len(input))
	for index := range input {
		output[index] = routeDto{
			DomainName: &input[index].DomainName,
			Backends:   toBackendDtos(input[index].Backends),
		}
	}

	return &output
}

func toBackendDtos(input []stream.Backend) *[]backendDto {
	if input == nil {
		return nil
	}

	output := make([]backendDto, len(input))
	for index := range input {
		output[index] = *toBackendDto(&input[index])
	}

	return &output
}

func toRoutes(input *[]routeDto) []stream.Route {
	if input == nil {
		return nil
	}

	output := make([]stream.Route, len(*input))
	for index := range *input {
		output[index] = stream.Route{
			DomainName: getStringValue((*input)[index].DomainName),
			Backends:   toBackends((*input)[index].Backends),
		}
	}

	return output
}

func toBackends(input *[]backendDto) []stream.Backend {
	if input == nil {
		return nil
	}

	output := make([]stream.Backend, len(*input))
	for index := range *input {
		output[index] = *toBackend(&(*input)[index])
	}

	return output
}

func toFeatureSet(input *featureSetDto) *stream.FeatureSet {
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

func toAddress(input *addressDto) *stream.Address {
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
