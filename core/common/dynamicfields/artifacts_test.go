package dynamicfields

import (
	"context"

	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

func newDynamicField(ctx context.Context) *DynamicField {
	return &DynamicField{
		ID:          "field1",
		Description: i18n.M(ctx, i18n.K.CertificateCustomName),
		Type:        SingleLineTextType,
		Priority:    100,
		Required:    false,
		Sensitive:   false,
	}
}
