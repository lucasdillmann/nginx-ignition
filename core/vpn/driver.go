package vpn

import (
	"context"

	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

type Destination struct {
	DomainName string
	IP         string
	Port       int
	HTTPS      bool
}

type Driver interface {
	ID() string
	Name() string
	ImportantInstructions() []string
	ConfigurationFields() []*dynamic_fields.DynamicField
	Reload(ctx context.Context, name, configDir string, destination *Destination, parameters map[string]any) error
	Start(ctx context.Context, name, configDir string, destination *Destination, parameters map[string]any) error
	Stop(ctx context.Context, name string) error
}
