package stream

import (
	"context"
	"strconv"
	"strings"

	"dillmann.com.br/nginx-ignition/core/common/constants"
	"dillmann.com.br/nginx-ignition/core/common/validation"
)

type validator struct {
	delegate *validation.ConsistencyValidator
}

func newValidator() *validator {
	return &validator{
		delegate: validation.NewValidator(),
	}
}

const (
	minimumPort  = 1
	maximumPort  = 65535
	invalidValue = "Invalid value"
)

func (v *validator) validate(_ context.Context, stream *Stream) error {
	if stream == nil {
		v.delegate.Add("", "Stream cannot be nil")
		return v.delegate.Result()
	}

	v.validateName(stream)
	v.validateType(stream)
	v.validateBinding(stream)
	v.validateDefaultBackend(stream)
	v.validateRoutes(stream)
	v.validateFeatureSet(stream)

	return v.delegate.Result()
}

func (v *validator) validateName(stream *Stream) {
	if strings.TrimSpace(stream.Name) == "" {
		v.delegate.Add("name", "Name cannot be empty")
	}
}

func (v *validator) validateType(stream *Stream) {
	switch stream.Type {
	case SimpleType, SNIRouterType:
	default:
		v.delegate.Add("type", invalidValue)
	}
}

func (v *validator) validateBinding(stream *Stream) {
	v.validateAddress("binding", stream.Binding)
}

func (v *validator) validateDefaultBackend(stream *Stream) {
	v.validateAddress("defaultBackend.target", stream.DefaultBackend.Address)
	v.validateCircuitBreaker("defaultBackend.circuitBreaker", stream.DefaultBackend.CircuitBreaker)
}

func (v *validator) validateRoutes(stream *Stream) {
	if stream.Type != SNIRouterType {
		return
	}

	if len(stream.Routes) == 0 {
		v.delegate.Add("routes", "Must be informed and not be empty when type is SNI_ROUTER")
		return
	}

	for index := range stream.Routes {
		v.validateRoute(&stream.Routes[index], index)
	}
}

func (v *validator) validateRoute(route *Route, index int) {
	prefix := "routes[" + strconv.Itoa(index) + "]"

	if len(route.DomainNames) == 0 {
		v.delegate.Add(prefix+".domainNames", "Route must have at least one domain")
	} else {
		for domainNameIndex, domainName := range route.DomainNames {
			v.validateDomainName(domainName, prefix, domainNameIndex)
		}
	}

	if len(route.Backends) == 0 {
		v.delegate.Add(prefix+".backends", "Route must have at least one backend")
	} else {
		for backendIndex, backend := range route.Backends {
			v.validateBackend(&backend, prefix, backendIndex)
		}
	}
}

func (v *validator) validateDomainName(domain, prefix string, index int) {
	domainPrefix := prefix + ".domainNames[" + strconv.Itoa(index) + "]"

	if domain == "" {
		v.delegate.Add(domainPrefix, "Domain cannot be empty")
	} else if !constants.TLDPattern.MatchString(domain) {
		v.delegate.Add(domainPrefix, "Not a valid DNS domain name")
	}
}

func (v *validator) validateBackend(backend *Backend, routePrefix string, index int) {
	prefix := routePrefix + ".backends[" + strconv.Itoa(index) + "]"

	v.validateAddress(prefix+".target", backend.Address)
	v.validateCircuitBreaker(prefix+".circuitBreaker", backend.CircuitBreaker)
}

func (v *validator) validateCircuitBreaker(prefix string, circuitBreaker *CircuitBreaker) {
	if circuitBreaker == nil {
		return
	}

	if circuitBreaker.MaxFailures < 1 {
		v.delegate.Add(prefix+".maxFailures", "Value must be greater than or equal to 1")
	}

	if circuitBreaker.OpenSeconds < 0 {
		v.delegate.Add(prefix+".openSeconds", "Value must be greater than or equal to 0")
	}
}

func (v *validator) validateAddress(fieldPrefix string, address Address) {
	switch address.Protocol {
	case UDPProtocol, TCPProtocol, SocketProtocol:
		break
	default:
		v.delegate.Add(fieldPrefix+".protocol", invalidValue)
	}

	v.validateAddressValue(fieldPrefix, address)
	v.validateAddressProtocol(fieldPrefix, address)
}

func (v *validator) validateAddressProtocol(fieldPrefix string, address Address) {
	if address.Protocol != SocketProtocol {
		if address.Port == nil {
			v.delegate.Add(fieldPrefix+".port", "Port is required when using TCP or UDP protocol")
		} else if *address.Port < minimumPort || *address.Port > maximumPort {
			v.delegate.Add(
				fieldPrefix+".port",
				"Value must be between "+strconv.Itoa(minimumPort)+" and "+strconv.Itoa(maximumPort),
			)
		}
	} else if address.Port != nil {
		v.delegate.Add(fieldPrefix+".port", "Port should not be specified when using the Socket protocol")
	}
}

func (v *validator) validateAddressValue(fieldPrefix string, address Address) {
	if strings.TrimSpace(address.Address) == "" {
		v.delegate.Add(fieldPrefix+".address", "Address cannot be empty")
		return
	}

	if address.Protocol == SocketProtocol && !strings.HasPrefix(address.Address, "/") {
		v.delegate.Add(fieldPrefix+".protocol", "Unix socket path must start with a /")
		return
	}
}

func (v *validator) validateFeatureSet(stream *Stream) {
	if stream.Binding.Protocol == TCPProtocol {
		return
	}

	if stream.FeatureSet.TCPKeepAlive {
		v.delegate.Add("featureSet.tcpKeepAlive", "TCP Keep Alive can be enabled only when binding uses the TCP protocol")
	}

	if stream.FeatureSet.TCPNoDelay {
		v.delegate.Add("featureSet.tcpNoDelay", "TCP No Delay can be enabled only when binding uses the TCP protocol")
	}

	if stream.FeatureSet.TCPDeferred {
		v.delegate.Add("featureSet.tcpDeferred", "TCP Deferred can be enabled only when binding uses the TCP protocol")
	}
}
