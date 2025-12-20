package integration

import (
	"context"
	"sort"

	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/common/coreerror"
	"dillmann.com.br/nginx-ignition/core/common/pagination"
)

type service struct {
	repository Repository
	drivers    func() []Driver
}

func newService(repository Repository, drivers func() []Driver) *service {
	return &service{
		repository: repository,
		drivers:    drivers,
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

func (s *service) getByID(ctx context.Context, id uuid.UUID) (*Integration, error) {
	return s.repository.FindByID(ctx, id)
}

func (s *service) save(ctx context.Context, data *Integration) error {
	driver := s.findDriver(data)
	if err := newValidator(s.repository, driver).validate(ctx, data); err != nil {
		return err
	}

	return s.repository.Save(ctx, data)
}

func (s *service) deleteByID(ctx context.Context, id uuid.UUID) error {
	inUse, err := s.repository.InUseByID(ctx, id)
	if err != nil {
		return err
	}

	if *inUse {
		return coreerror.New("Integration is in use by one or more hosts", true)
	}

	return s.repository.DeleteByID(ctx, id)
}

func (s *service) existsByID(ctx context.Context, id uuid.UUID) (*bool, error) {
	return s.repository.ExistsByID(ctx, id)
}

func (s *service) listOptions(
	ctx context.Context,
	integrationId uuid.UUID,
	pageNumber, pageSize int,
	searchTerms *string,
	tcpOnly bool,
) (*pagination.Page[*DriverOption], error) {
	data, err := s.repository.FindByID(ctx, integrationId)
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, ErrIntegrationNotFound
	}

	if !data.Enabled {
		return nil, ErrIntegrationDisabled
	}

	driver := s.findDriver(data)
	if driver == nil {
		return nil, ErrIntegrationNotFound
	}

	options, err := driver.GetAvailableOptions(ctx, data.Parameters, pageNumber, pageSize, searchTerms, tcpOnly)
	if err != nil {
		return nil, err
	}

	sort.Slice(options.Contents, func(left, right int) bool {
		return options.Contents[left].Name < options.Contents[right].Name
	})

	return options, nil
}

func (s *service) getOptionByID(ctx context.Context, integrationId uuid.UUID, optionId string) (*DriverOption, error) {
	data, err := s.repository.FindByID(ctx, integrationId)
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, ErrIntegrationNotFound
	}

	driver := s.findDriver(data)
	if driver == nil {
		return nil, ErrIntegrationNotFound
	}

	if !data.Enabled {
		return nil, ErrIntegrationDisabled
	}

	return driver.GetAvailableOptionById(ctx, data.Parameters, optionId)
}

func (s *service) getOptionURL(
	ctx context.Context,
	integrationId uuid.UUID,
	optionId string,
) (*string, *[]string, error) {
	data, err := s.repository.FindByID(ctx, integrationId)
	if err != nil {
		return nil, nil, err
	}

	if data == nil {
		return nil, nil, ErrIntegrationNotFound
	}

	driver := s.findDriver(data)
	if driver == nil {
		return nil, nil, ErrIntegrationNotFound
	}

	if !data.Enabled {
		return nil, nil, ErrIntegrationDisabled
	}

	return driver.GetOptionProxyURL(ctx, data.Parameters, optionId)
}

func (s *service) getAvailableDrivers(_ context.Context) (*[]*AvailableDriver, error) {
	drivers := s.drivers()
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
	for _, driver := range s.drivers() {
		if driver.ID() == data.Driver {
			return driver
		}
	}

	return nil
}
