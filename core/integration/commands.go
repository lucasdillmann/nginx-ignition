package integration

import (
	"context"

	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/pagination"
)

type AvailableDriver struct {
	ID                  string
	Name                string
	Description         string
	ConfigurationFields []*dynamicfields.DynamicField
}

type Commands struct {
	Get                 func(ctx context.Context, id uuid.UUID) (*Integration, error)
	Delete              func(ctx context.Context, id uuid.UUID) error
	Save                func(ctx context.Context, data *Integration) error
	Exists              func(ctx context.Context, id uuid.UUID) (*bool, error)
	GetOption           func(ctx context.Context, integrationId uuid.UUID, optionId string) (*DriverOption, error)
	GetOptionURL        func(ctx context.Context, integrationId uuid.UUID, optionId string) (*string, *[]string, error)
	GetAvailableDrivers func(ctx context.Context) (*[]*AvailableDriver, error)
	List                func(
		ctx context.Context,
		pageSize, pageNumber int,
		searchTerms *string,
		enabledOnly bool,
	) (*pagination.Page[*Integration], error)
	ListOptions func(
		ctx context.Context,
		integrationId uuid.UUID,
		pageNumber, pageSize int,
		searchTerms *string,
		tcpOnly bool,
	) (*pagination.Page[*DriverOption], error)
}
