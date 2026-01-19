package stream

import (
	"context"
	"fmt"
	"strings"

	"dillmann.com.br/nginx-ignition/core/common/constants"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
	"dillmann.com.br/nginx-ignition/core/common/validation"
	"dillmann.com.br/nginx-ignition/core/common/valuerange"
)

var portRange = valuerange.New(1, 65535)

type validator struct {
	delegate *validation.ConsistencyValidator
}

func newValidator() *validator {
	return &validator{
		delegate: validation.NewValidator(),
	}
}

func (v *validator) validate(ctx context.Context, stream *Stream) error {
	if stream == nil {
		v.delegate.Add("", i18n.M(ctx, i18n.K.CoreStreamNilStream))
		return v.delegate.Result()
	}

	v.validateName(ctx, stream)
	v.validateType(ctx, stream)
	v.validateBinding(ctx, stream)
	v.validateDefaultBackend(ctx, stream)
	v.validateRoutes(ctx, stream)
	v.validateFeatureSet(ctx, stream)

	return v.delegate.Result()
}

func (v *validator) validateName(ctx context.Context, stream *Stream) {
	if strings.TrimSpace(stream.Name) == "" {
		v.delegate.Add("name", i18n.M(ctx, i18n.K.CommonCannotBeEmpty))
	}
}

func (v *validator) validateType(ctx context.Context, stream *Stream) {
	switch stream.Type {
	case SimpleType, SNIRouterType:
	default:
		v.delegate.Add("type", i18n.M(ctx, i18n.K.CommonInvalidValue))
	}
}

func (v *validator) validateBinding(ctx context.Context, stream *Stream) {
	v.validateAddress(ctx, "binding", stream.Binding)
}

func (v *validator) validateDefaultBackend(ctx context.Context, stream *Stream) {
	v.validateAddress(ctx, "defaultBackend.target", stream.DefaultBackend.Address)
	v.validateCircuitBreaker(
		ctx,
		"defaultBackend.circuitBreaker",
		stream.DefaultBackend.CircuitBreaker,
	)
}

func (v *validator) validateRoutes(ctx context.Context, stream *Stream) {
	if stream.Type != SNIRouterType {
		return
	}

	if len(stream.Routes) == 0 {
		v.delegate.Add("routes", i18n.M(ctx, i18n.K.CoreStreamRoutesRequiredForSni))
		return
	}

	for index := range stream.Routes {
		v.validateRoute(ctx, &stream.Routes[index], index)
	}
}

func (v *validator) validateRoute(ctx context.Context, route *Route, index int) {
	prefix := fmt.Sprintf("routes[%d]", index)

	if len(route.DomainNames) == 0 {
		v.delegate.Add(prefix+".domainNames", i18n.M(ctx, i18n.K.CoreStreamAtLeastOneDomain))
	} else {
		for domainNameIndex, domainName := range route.DomainNames {
			v.validateDomainName(ctx, domainName, prefix, domainNameIndex)
		}
	}

	if len(route.Backends) == 0 {
		v.delegate.Add(prefix+".backends", i18n.M(ctx, i18n.K.CoreStreamAtLeastOneBackend))
	} else {
		for backendIndex, backend := range route.Backends {
			v.validateBackend(ctx, &backend, prefix, backendIndex)
		}
	}
}

func (v *validator) validateDomainName(ctx context.Context, domain, prefix string, index int) {
	domainPrefix := fmt.Sprintf("%s.domainNames[%d]", prefix, index)

	if domain == "" {
		v.delegate.Add(domainPrefix, i18n.M(ctx, i18n.K.CommonCannotBeEmpty))
	} else if !constants.TLDPattern.MatchString(domain) {
		v.delegate.Add(domainPrefix, i18n.M(ctx, i18n.K.CommonInvalidDomainName))
	}
}

func (v *validator) validateBackend(
	ctx context.Context,
	backend *Backend,
	routePrefix string,
	index int,
) {
	prefix := fmt.Sprintf("%s.backends[%d]", routePrefix, index)

	v.validateAddress(ctx, prefix+".target", backend.Address)
	v.validateCircuitBreaker(ctx, prefix+".circuitBreaker", backend.CircuitBreaker)
}

func (v *validator) validateCircuitBreaker(
	ctx context.Context,
	prefix string,
	circuitBreaker *CircuitBreaker,
) {
	if circuitBreaker == nil {
		return
	}

	if circuitBreaker.MaxFailures < 1 {
		v.delegate.Add(prefix+".maxFailures", i18n.M(ctx, i18n.K.CommonCannotBeZero))
	}

	if circuitBreaker.OpenSeconds < 0 {
		v.delegate.Add(prefix+".openSeconds", i18n.M(ctx, i18n.K.CoreStreamCannotBeNegative))
	}
}

func (v *validator) validateAddress(ctx context.Context, fieldPrefix string, address Address) {
	switch address.Protocol {
	case UDPProtocol, TCPProtocol, SocketProtocol:
	default:
		v.delegate.Add(fieldPrefix+".protocol", i18n.M(ctx, i18n.K.CommonInvalidValue))
	}

	v.validateAddressValue(ctx, fieldPrefix, address)
	v.validateAddressProtocol(ctx, fieldPrefix, address)
}

func (v *validator) validateAddressProtocol(
	ctx context.Context,
	fieldPrefix string,
	address Address,
) {
	if address.Protocol != SocketProtocol {
		if address.Port == nil {
			v.delegate.Add(
				fieldPrefix+".port",
				i18n.M(ctx, i18n.K.CoreStreamPortRequired),
			)
		} else if !portRange.Contains(*address.Port) {
			v.delegate.Add(
				fieldPrefix+".port",
				i18n.M(ctx, i18n.K.CommonBetweenValues).
					V("min", portRange.Min).
					V("max", portRange.Max),
			)
		}
	} else if address.Port != nil {
		v.delegate.Add(
			fieldPrefix+".port",
			i18n.M(ctx, i18n.K.CoreStreamPortNotAllowedForSocket),
		)
	}
}

func (v *validator) validateAddressValue(ctx context.Context, fieldPrefix string, address Address) {
	if strings.TrimSpace(address.Address) == "" {
		v.delegate.Add(fieldPrefix+".address", i18n.M(ctx, i18n.K.CommonCannotBeEmpty))
		return
	}

	if address.Protocol == SocketProtocol && !strings.HasPrefix(address.Address, "/") {
		v.delegate.Add(
			fieldPrefix+".protocol",
			i18n.M(ctx, i18n.K.CommonStartsWithSlashRequired),
		)
		return
	}
}

func (v *validator) validateFeatureSet(ctx context.Context, stream *Stream) {
	if stream.Binding.Protocol == TCPProtocol {
		return
	}

	if stream.FeatureSet.TCPKeepAlive {
		v.delegate.Add(
			"featureSet.tcpKeepAlive",
			i18n.M(ctx, i18n.K.CoreStreamFeatureOnlyForTcp).V("feature", "TCP Keep Alive"),
		)
	}

	if stream.FeatureSet.TCPNoDelay {
		v.delegate.Add(
			"featureSet.tcpNoDelay",
			i18n.M(ctx, i18n.K.CoreStreamFeatureOnlyForTcp).V("feature", "TCP No Delay"),
		)
	}

	if stream.FeatureSet.TCPDeferred {
		v.delegate.Add(
			"featureSet.tcpDeferred",
			i18n.M(ctx, i18n.K.CoreStreamFeatureOnlyForTcp).V("feature", "TCP Deferred"),
		)
	}
}
