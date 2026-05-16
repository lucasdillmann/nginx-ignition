package external

import (
	"context"

	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

const (
	publicKeyPathFieldID  = "publicKeyPath"
	privateKeyPathFieldID = "privateKeyPath"
	chainPathFieldID      = "chainPath"
)

func dynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	publicKeyPathField := dynamicfields.DynamicField{
		ID:          publicKeyPathFieldID,
		Priority:    0,
		Description: i18n.M(ctx, i18n.K.CertificateExternalPublicKeyPath),
		Required:    true,
		Sensitive:   false,
		Type:        dynamicfields.SingleLineTextType,
	}

	privateKeyPathField := dynamicfields.DynamicField{
		ID:          privateKeyPathFieldID,
		Priority:    1,
		Description: i18n.M(ctx, i18n.K.CertificateExternalPrivateKeyPath),
		Required:    true,
		Sensitive:   false,
		Type:        dynamicfields.SingleLineTextType,
	}

	chainPathField := dynamicfields.DynamicField{
		ID:          chainPathFieldID,
		Priority:    2,
		Description: i18n.M(ctx, i18n.K.CertificateExternalChainPath),
		Required:    false,
		Sensitive:   false,
		Type:        dynamicfields.SingleLineTextType,
	}

	return []dynamicfields.DynamicField{
		publicKeyPathField,
		privateKeyPathField,
		chainPathField,
	}
}
