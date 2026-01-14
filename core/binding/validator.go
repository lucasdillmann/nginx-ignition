package binding

import (
	"context"
	"fmt"
	"net"

	"dillmann.com.br/nginx-ignition/core/certificate"
	"dillmann.com.br/nginx-ignition/core/common/validation"
	"dillmann.com.br/nginx-ignition/core/common/valuerange"
	"dillmann.com.br/nginx-ignition/i18n"
)

var portRange = valuerange.New(1, 65535)

type validator struct {
	delegate            *validation.ConsistencyValidator
	certificateCommands certificate.Commands
}

func newValidator(
	validationCtx *validation.ConsistencyValidator,
	certificateCommands certificate.Commands,
) *validator {
	return &validator{
		delegate:            validationCtx,
		certificateCommands: certificateCommands,
	}
}

func (v *validator) validate(
	ctx context.Context,
	pathPrefix string,
	binding *Binding,
	index int,
) error {
	if net.ParseIP(binding.IP) == nil {
		v.delegate.Add(
			fmt.Sprintf("%s[%d].ip", pathPrefix, index),
			i18n.M(ctx, i18n.K.BindingValidationInvalidIP),
		)
	}

	if !portRange.Contains(binding.Port) {
		v.delegate.Add(
			fmt.Sprintf("%s[%d].port", pathPrefix, index),
			i18n.M(ctx, i18n.K.CommonValidationBetweenValues).
				V("min", portRange.Min).
				V("max", portRange.Max),
		)
	}

	certificateIDField := fmt.Sprintf("%s[%d].certificateId", pathPrefix, index)

	switch {
	case binding.Type == HTTPBindingType && binding.CertificateID != nil:
		v.delegate.Add(
			certificateIDField,
			i18n.M(ctx, i18n.K.BindingValidationCertificateIDNotAllowed),
		)
	case binding.Type == HTTPBindingType && binding.CertificateID == nil:
		return nil
	case binding.Type == HTTPSBindingType && binding.CertificateID == nil:
		v.delegate.Add(
			certificateIDField,
			i18n.M(ctx, i18n.K.BindingValidationCertificateIDRequired),
		)
	case binding.Type == HTTPSBindingType:
		exists, err := v.certificateCommands.Exists(ctx, *binding.CertificateID)
		if err != nil {
			return err
		}

		if !exists {
			v.delegate.Add(
				certificateIDField,
				i18n.M(ctx, i18n.K.BindingValidationCertificateIDNotFound),
			)
		}
	default:
		v.delegate.Add(
			fmt.Sprintf("%s[%d].certificateId", pathPrefix, index),
			i18n.M(ctx, i18n.K.BindingValidationInvalidType),
		)
	}

	return nil
}
