package vpn

import (
	"context"

	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
	"dillmann.com.br/nginx-ignition/core/host"
)

type Driver interface {
	ID() string
	Name() string
	Description() string
	ConfigurationFields() []*dynamic_fields.DynamicField
	Start(ctx context.Context, h *host.Host, parameters map[string]any) error
	Stop(ctx context.Context, h *host.Host, parameters map[string]any) error
}
