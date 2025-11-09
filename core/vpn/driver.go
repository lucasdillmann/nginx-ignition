package vpn

import (
	"context"

	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

type Destination interface {
	VPNID() uuid.UUID
	Name() string
	DomainName() string
	IP() string
	Port() int
	HTTPS() bool
}

type Driver interface {
	ID() string
	Name() string
	ImportantInstructions() []string
	ConfigurationFields() []*dynamic_fields.DynamicField
	Reload(ctx context.Context, configDir string, destination Destination, parameters map[string]any) error
	Start(ctx context.Context, configDir string, destination Destination, parameters map[string]any) error
	Stop(ctx context.Context, destination Destination) error
}
