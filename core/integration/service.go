package integration

import (
	"context"
	"sort"

	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/common/core_error"
	"dillmann.com.br/nginx-ignition/core/common/pagination"
)

type service struct {
	repository      Repository
	driversResolver func() ([]Driver, error)
}

func newService(repository Repository, adaptersResolver func() ([]Driver, error)) *service {
	return &service{
		repository:      repository,
		driversResolver: adaptersResolver,
	}
}

func (s *service) list(
	ctx context.Context,
	pageSize, pageNumber int,
	searchTerms *string,
	enabledOnly bool,
) (*pagination.Page[*Integration], error) {
	return s.repository.FindPage(ctx, pageSize, pageNumber, searchTerms, enabledOnly)
}

func (s *service) getById(ctx context.Context, id uuid.UUID) (*Integration, error) {
	return s.repository.FindById(ctx, id)
}

func (s *service) save(ctx context.Context, data *Integration) error {
	driver := s.findDriver(data)
	if err := newValidator(s.repository, driver).validate(ctx, data); err != nil {
		return err
	}

	return s.repository.Save(ctx, data)
}

func (s *service) deleteById(ctx context.Context, id uuid.UUID) error {
	inUse, err := s.repository.ExistsByID(ctx, id)
	if err != nil {
		return err
	}

	if *inUse {
		return core_error.New("Integration is in use by one or more hosts", true)
	}

	return s.repository.DeleteById(ctx, id)
}

func (s *service) existsById(ctx context.Context, id uuid.UUID) (*bool, error) {
	return s.repository.ExistsByID(ctx, id)
}

func (s *service) listOptions(
	ctx context.Context,
	integrationId uuid.UUID,
	pageNumber, pageSize int,
	searchTerms *string,
	tcpOnly bool,
) (*pagination.Page[*DriverOption], error) {
	data, err := s.repository.FindById(ctx, integrationId)
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, ErrIntegrationNotFound
	}

	if !data.Enabled {
		return nil, ErrIntegrationDisabled
	}

	adapter := s.findDriver(data)
	if adapter == nil {
		return nil, ErrIntegrationNotFound
	}

	options, err := adapter.GetAvailableOptions(ctx, data.Parameters, pageNumber, pageSize, searchTerms, tcpOnly)
	if err != nil {
		return nil, err
	}

	sort.Slice(options.Contents, func(i, j int) bool {
		return options.Contents[i].Name < options.Contents[j].Name
	})

	return options, nil
}

func (s *service) getOptionById(ctx context.Context, integrationId uuid.UUID, optionId string) (*DriverOption, error) {
	data, err := s.repository.FindById(ctx, integrationId)
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, ErrIntegrationNotFound
	}

	adapter := s.findDriver(data)
	if adapter == nil {
		return nil, ErrIntegrationNotFound
	}

	if !data.Enabled {
		return nil, ErrIntegrationDisabled
	}

	return adapter.GetAvailableOptionById(ctx, data.Parameters, optionId)
}

func (s *service) getOptionUrl(ctx context.Context, integrationId uuid.UUID, optionId string) (*string, error) {
	data, err := s.repository.FindById(ctx, integrationId)
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, ErrIntegrationNotFound
	}

	adapter := s.findDriver(data)
	if adapter == nil {
		return nil, ErrIntegrationNotFound
	}

	if !data.Enabled {
		return nil, ErrIntegrationDisabled
	}

	url, err := adapter.GetOptionProxyURL(ctx, data.Parameters, optionId)
	if err != nil {
		return nil, err
	}

	return url, nil
}

func (s *service) getAvailableDrivers(_ context.Context) (*[]*AvailableDriver, error) {
	drivers, err := s.driversResolver()
	if err != nil {
		return nil, err
	}

	sort.Slice(drivers, func(left, right int) bool {
		return drivers[left].Name() < drivers[right].Name()
	})

	output := make([]*AvailableDriver, len(drivers))
	for index, driver := range drivers {
		output[index] = &AvailableDriver{
			ID:                  driver.ID(),
			Name:                driver.Name(),
			Description:         driver.Description(),
			ConfigurationFields: driver.ConfigurationFields(),
		}
	}

	return &output, nil
}

func (s *service) findDriver(data *Integration) Driver {
	drivers, err := s.driversResolver()
	if err != nil {
		return nil
	}

	for _, driver := range drivers {
		if driver.ID() == data.Driver {
			return driver
		}
	}

	return nil
}
