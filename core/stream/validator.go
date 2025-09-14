package stream

import (
	"context"
	"net"
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
	v.validateBinding(stream)
	v.validateBackend(stream)
	v.validateFeatureSet(stream)

	return v.delegate.Result()
}

func (v *validator) validateName(stream *Stream) {
	if strings.TrimSpace(stream.Name) == "" {
		v.delegate.Add("name", "Name cannot be empty")
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
	path := fieldPrefix + ".address"

	if strings.TrimSpace(address.Address) == "" {
		v.delegate.Add(path, "Address cannot be empty")
		return
	}

	if address.Protocol == SocketProtocol && !strings.HasPrefix(address.Address, "/") {
		v.delegate.Add(path, "Unix socket path must start with a /")
		return
	}

	if address.Protocol != SocketProtocol &&
		net.ParseIP(address.Address) == nil &&
		!constants.TLDPattern.MatchString(address.Address) {
		v.delegate.Add(path, "Not a valid IP address or domain name")
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
