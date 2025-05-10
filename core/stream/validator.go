package stream

import (
	"context"
	"dillmann.com.br/nginx-ignition/core/common/constants"
	"dillmann.com.br/nginx-ignition/core/common/validation"
	"net"
	"strconv"
	"strings"
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

	v.validateDescription(stream)
	v.validateBinding(stream)
	v.validateBackend(stream)
	v.validateFeatureSet(stream)

	return v.delegate.Result()
}

func (v *validator) validateDescription(stream *Stream) {
	if strings.TrimSpace(stream.Description) == "" {
		v.delegate.Add("description", "Description cannot be empty")
	}
}

func (v *validator) validateBinding(stream *Stream) {
	v.validateAddress("binding", stream.Binding)
}

func (v *validator) validateBackend(stream *Stream) {
	v.validateAddress("backend", stream.Backend)
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
			v.delegate.Add(fieldPrefix+".port", "Port is required when binding is using TCP or UDP protocol")
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
	} else {
		if net.ParseIP(address.Address) == nil {
			if address.Protocol != SocketProtocol && !constants.TLDPattern.MatchString(address.Address) {
				v.delegate.Add(fieldPrefix+".address", "Not a valid IP address or domain name")
			}
		}
	}
}

func (v *validator) validateFeatureSet(stream *Stream) {
	if stream.Binding.Protocol != UDPProtocol {
		return
	}

	if stream.FeatureSet.TCPKeepAlive {
		v.delegate.Add("featureSet.tcpKeepAlive", "TCP Keep Alive cannot be enabled when binding uses the UDP protocol")
	}

	if stream.FeatureSet.TCPNoDelay {
		v.delegate.Add("featureSet.tcpNoDelay", "TCP No Delay cannot be enabled when binding uses the UDP protocol")
	}

	if stream.FeatureSet.TCPDeferred {
		v.delegate.Add("featureSet.tcpDeferred", "TCP Deferred cannot be enabled when binding uses the UDP protocol")
	}

	if stream.FeatureSet.SSL {
		v.delegate.Add("featureSet.ssl", "SSL cannot be enabled when binding uses the UDP protocol")
	}

}
