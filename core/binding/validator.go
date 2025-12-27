package binding

import (
	"context"
	"fmt"
	"net"

	"dillmann.com.br/nginx-ignition/core/certificate"
	"dillmann.com.br/nginx-ignition/core/common/validation"
	"dillmann.com.br/nginx-ignition/core/common/valuerange"
)

var portRange = valuerange.New(1, 65535)

type validator struct {
	delegate            *validation.ConsistencyValidator
	certificateCommands *certificate.Commands
}

func newValidator(validationCtx *validation.ConsistencyValidator, certificateCommands *certificate.Commands) *validator {
	return &validator{
		delegate:            validationCtx,
		certificateCommands: certificateCommands,
	}
}

func (v *validator) validate(ctx context.Context, pathPrefix string, binding *Binding, index int) error {
	if net.ParseIP(binding.IP) == nil {
		v.delegate.Add(
			fmt.Sprintf("%s[%d].ip", pathPrefix, index),
			"Not a valid IPv4 or IPv6 address",
		)
	}

	if !portRange.Contains(binding.Port) {
		v.delegate.Add(
			fmt.Sprintf("%s[%d].port", pathPrefix, index),
			fmt.Sprintf("Value must be between %d and %d", portRange.Min, portRange.Max),
		)
	}

	certificateIdField := fmt.Sprintf("%s[%d].certificateId", pathPrefix, index)

	switch {
	case binding.Type == HttpBindingType && binding.CertificateID != nil:
		v.delegate.Add(certificateIdField, "Value cannot be informed for a HTTP binding")
	case binding.Type == HttpBindingType && binding.CertificateID == nil:
		return nil
	case binding.Type == HttpsBindingType && binding.CertificateID == nil:
		v.delegate.Add(certificateIdField, "Value must be informed for a HTTPS binding")
	case binding.Type == HttpsBindingType:
		exists, err := v.certificateCommands.Exists(ctx, *binding.CertificateID)
		if err != nil {
			return err
		}

		if !exists {
			v.delegate.Add(certificateIdField, "No SSL certificate found with provided ID")
		}
	default:
		v.delegate.Add(
			fmt.Sprintf("%s[%d].certificateId", pathPrefix, index),
			"Invalid binding type",
		)
	}

	return nil
}
