package vpn

import (
	"context"

	"github.com/google/uuid"

	"nginx-ignition/internal/core/common/dynamicfields"
	"nginx-ignition/internal/core/common/i18n"
	"nginx-ignition/internal/core/common/pagination"
)

type AvailableDriver struct {
	Name                  *i18n.Message
	ID                    string
	EndpointSSLSupport    EndpointSSLSupport
	ImportantInstructions []*i18n.Message
	ConfigurationFields   []dynamicfields.DynamicField
}

type Commands interface {
	Get(ctx context.Context, id uuid.UUID) (*VPN, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Save(ctx context.Context, data *VPN) error
	Exists(ctx context.Context, id uuid.UUID) (*bool, error)
	GetAvailableDrivers(ctx context.Context) ([]AvailableDriver, error)
	GetAvailableDriverByID(ctx context.Context, id string) (*AvailableDriver, error)
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
