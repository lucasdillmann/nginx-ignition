package vpn

import (
	"context"

	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

type Driver interface {
	ID() string
	Name() string
	ImportantInstructions(ctx context.Context) []*i18n.Message
	ConfigurationFields(ctx context.Context) []dynamicfields.DynamicField
	Reload(
		ctx context.Context,
		configDir string,
		endpoint Endpoint,
		parameters map[string]any,
	) error
	Start(ctx context.Context, configDir string, endpoint Endpoint, parameters map[string]any) error
	Stop(ctx context.Context, endpoint Endpoint) error
}
