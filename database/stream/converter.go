package stream

import (
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/stream"
)

func toDomain(model *streamModel) *stream.Stream {
	if model == nil {
		return nil
	}

	return &stream.Stream{
		ID:             model.ID,
		Enabled:        model.Enabled,
		Name:           model.Name,
		Type:           stream.Type(model.Type),
		DefaultBackend: stream.Backend{},
		Binding: stream.Address{
			Protocol: stream.Protocol(model.BindingProtocol),
			Address:  model.BindingAddress,
			Port:     model.BindingPort,
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

func toDomainBackend(model *streamBackendModel) *stream.Backend {
	if model == nil {
		return nil
	}

	circuitBreaker := &stream.CircuitBreaker{}
	if model.MaxFailures == nil || model.OpenSeconds == nil {
		circuitBreaker = nil
	} else {
		circuitBreaker.MaxFailures = *model.MaxFailures
		circuitBreaker.OpenSeconds = *model.OpenSeconds
	}

	return &stream.Backend{
		Weight:         model.Weight,
		CircuitBreaker: circuitBreaker,
		Address: stream.Address{
			Protocol: stream.Protocol(model.Protocol),
			Address:  model.Address,
			Port:     model.Port,
		},
	}
}

func toDomainRoute(model *streamRouteModel, backendModels []streamBackendModel) *stream.Route {
	if model == nil {
		return nil
	}

	backends := make([]stream.Backend, len(backendModels))
	for index, backend := range backendModels {
		backends[index] = *toDomainBackend(&backend)
	}

	return &stream.Route{
		DomainNames: model.DomainNames,
		Backends:    backends,
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
		Type:             string(domain.Type),
		BindingProtocol:  string(domain.Binding.Protocol),
		BindingAddress:   domain.Binding.Address,
		BindingPort:      domain.Binding.Port,
		UseProxyProtocol: domain.FeatureSet.UseProxyProtocol,
		SocketKeepAlive:  domain.FeatureSet.SocketKeepAlive,
		TCPKeepAlive:     domain.FeatureSet.TCPKeepAlive,
		TCPNoDelay:       domain.FeatureSet.TCPNoDelay,
		TCPDeferred:      domain.FeatureSet.TCPDeferred,
	}
}

func toModelBackend(backend *stream.Backend, streamId, routeId *uuid.UUID) *streamBackendModel {
	if backend == nil {
		return nil
	}

	var maxFailures, openSeconds *int
	if backend.CircuitBreaker != nil {
		maxFailures = &backend.CircuitBreaker.MaxFailures
		openSeconds = &backend.CircuitBreaker.OpenSeconds
	}

	return &streamBackendModel{
		ID:            uuid.New(),
		StreamID:      streamId,
		StreamRouteID: routeId,
		Protocol:      string(backend.Address.Protocol),
		Address:       backend.Address.Address,
		Port:          backend.Address.Port,
		Weight:        backend.Weight,
		MaxFailures:   maxFailures,
		OpenSeconds:   openSeconds,
	}
}

func toModelRoute(route *stream.Route, streamId uuid.UUID) *streamRouteModel {
	if route == nil {
		return nil
	}

	return &streamRouteModel{
		ID:          uuid.New(),
		StreamID:    streamId,
		DomainNames: route.DomainNames,
	}
}
