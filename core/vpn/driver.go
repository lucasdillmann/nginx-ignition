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
	ConfigurationFields() []*dynamic_fields.DynamicField
	Start(ctx context.Context, name, configDir string, destination *Destination, parameters map[string]any) error
	Stop(ctx context.Context, name string) error
}
