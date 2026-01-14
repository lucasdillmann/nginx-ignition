package integration

import (
	"context"

	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
	"dillmann.com.br/nginx-ignition/core/common/pagination"
)

type AvailableDriver struct {
	ID                  string
	Name                string
	Description         *i18n.Message
	ConfigurationFields []dynamicfields.DynamicField
}

type Commands interface {
	Get(ctx context.Context, id uuid.UUID) (*Integration, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Save(ctx context.Context, data *Integration) error
	Exists(ctx context.Context, id uuid.UUID) (*bool, error)
	GetOption(ctx context.Context, integrationID uuid.UUID, optionID string) (*DriverOption, error)
	GetOptionURL(
		ctx context.Context,
		integrationID uuid.UUID,
		optionID string,
	) (*string, []string, error)
	GetAvailableDrivers(ctx context.Context) ([]AvailableDriver, error)
	List(
		ctx context.Context,
		pageSize, pageNumber int,
		searchTerms *string,
		enabledOnly bool,
	) (*pagination.Page[Integration], error)
	ListOptions(
		ctx context.Context,
		integrationID uuid.UUID,
		pageNumber, pageSize int,
		searchTerms *string,
		tcpOnly bool,
	) (*pagination.Page[DriverOption], error)
}
