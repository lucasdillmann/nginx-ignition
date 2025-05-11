package stream

import (
	"dillmann.com.br/nginx-ignition/core/stream"
)

func toDto(input *stream.Stream) *streamResponseDto {
	if input == nil {
		return nil
	}

	return &streamResponseDto{
		ID:         &input.ID,
		Enabled:    &input.Enabled,
		Name:       &input.Name,
		FeatureSet: toFeatureSetDto(&input.FeatureSet),
		Backend:    toAddressDto(&input.Backend),
		Binding:    toAddressDto(&input.Binding),
	}
}

func toDomain(input *streamRequestDto) *stream.Stream {
	if input == nil {
		return nil
	}

	return &stream.Stream{
		Enabled:    getBoolValue(input.Enabled),
		Name:       getStringValue(input.Name),
		FeatureSet: *toFeatureSet(input.FeatureSet),
		Backend:    *toAddress(input.Backend),
		Binding:    *toAddress(input.Binding),
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
