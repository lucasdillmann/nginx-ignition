package dynamicfields

import (
	"context"

	"github.com/lucasdillmann/nginx-ignition/internal/core/common/i18n"
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
