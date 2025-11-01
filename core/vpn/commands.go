package vpn

import (
	"context"

	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
	"dillmann.com.br/nginx-ignition/core/common/pagination"
)

type AvailableDriver struct {
	ID                  string
	Name                string
	ConfigurationFields []*dynamic_fields.DynamicField
}

type Commands struct {
	Get                 func(ctx context.Context, id uuid.UUID) (*VPN, error)
	Delete              func(ctx context.Context, id uuid.UUID) error
	Save                func(ctx context.Context, data *VPN) error
	Exists              func(ctx context.Context, id uuid.UUID) (*bool, error)
	GetAvailableDrivers func(ctx context.Context) (*[]*AvailableDriver, error)
	List                func(
		ctx context.Context,
		pageSize, pageNumber int,
		searchTerms *string,
		enabledOnly bool,
	) (*pagination.Page[*VPN], error)
}
