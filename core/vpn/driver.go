package vpn

import (
	"context"

	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

type Driver interface {
	ID() string
	Name() string
	ImportantInstructions() []string
	ConfigurationFields() []dynamicfields.DynamicField
	Reload(ctx context.Context, configDir string, endpoint Endpoint, parameters map[string]any) error
	Start(ctx context.Context, configDir string, endpoint Endpoint, parameters map[string]any) error
	Stop(ctx context.Context, endpoint Endpoint) error
}
