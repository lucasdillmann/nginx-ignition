package vpn

import (
	"context"

	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/pagination"
)

type AvailableDriver struct {
	ID                    string
	Name                  string
	ImportantInstructions []string
	ConfigurationFields   []dynamicfields.DynamicField
}

type Commands interface {
	Get(ctx context.Context, id uuid.UUID) (*VPN, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Save(ctx context.Context, data *VPN) error
	Exists(ctx context.Context, id uuid.UUID) (*bool, error)
	GetAvailableDrivers(ctx context.Context) ([]AvailableDriver, error)
	Start(ctx context.Context, endpoint Endpoint) error
	Reload(ctx context.Context, endpoint Endpoint) error
	Stop(ctx context.Context, endpoint Endpoint) error
	List(
		ctx context.Context,
		pageSize, pageNumber int,
		searchTerms *string,
		enabledOnly bool,
	) (*pagination.Page[VPN], error)
}
