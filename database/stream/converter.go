package stream

import (
	"dillmann.com.br/nginx-ignition/core/stream"
)

func toDomain(model *streamModel) *stream.Stream {
	if model == nil {
		return nil
	}

	return &stream.Stream{
		ID:      model.ID,
		Enabled: model.Enabled,
		Name:    model.Name,
		Binding: stream.Address{
			Protocol: stream.Protocol(model.BindingProtocol),
			Address:  model.BindingAddress,
			Port:     model.BindingPort,
		},
		Backend: stream.Address{
			Protocol: stream.Protocol(model.BackendProtocol),
			Address:  model.BackendAddress,
			Port:     model.BackendPort,
		},
		FeatureSet: stream.FeatureSet{
			UseProxyProtocol: model.UseProxyProtocol,
			SocketKeepAlive:  model.SocketKeepAlive,
			TCPKeepAlive:     model.TCPKeepAlive,
			TCPNoDelay:       model.TCPNoDelay,
			TCPDeferred:      model.TCPDeferred,
		},
	}
}

func toModel(domain *stream.Stream) *streamModel {
	if domain == nil {
		return nil
	}

	return &streamModel{
		ID:               domain.ID,
		Enabled:          domain.Enabled,
		Name:             domain.Name,
		BindingProtocol:  string(domain.Binding.Protocol),
		BindingAddress:   domain.Binding.Address,
		BindingPort:      domain.Binding.Port,
		BackendProtocol:  string(domain.Backend.Protocol),
		BackendAddress:   domain.Backend.Address,
		BackendPort:      domain.Backend.Port,
		UseProxyProtocol: domain.FeatureSet.UseProxyProtocol,
		SocketKeepAlive:  domain.FeatureSet.SocketKeepAlive,
		TCPKeepAlive:     domain.FeatureSet.TCPKeepAlive,
		TCPNoDelay:       domain.FeatureSet.TCPNoDelay,
		TCPDeferred:      domain.FeatureSet.TCPDeferred,
	}
}
