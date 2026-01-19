package integration

import (
	"context"
	"sort"

	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/common/coreerror"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
	"dillmann.com.br/nginx-ignition/core/common/pagination"
)

type service struct {
	repository Repository
	drivers    func() []Driver
}

func newService(repository Repository, drivers func() []Driver) Commands {
	return &service{
		repository: repository,
		drivers:    drivers,
	}
}

func (s *service) List(
	ctx context.Context,
	pageSize, pageNumber int,
	searchTerms *string,
	enabledOnly bool,
) (*pagination.Page[Integration], error) {
	return s.repository.FindPage(ctx, pageSize, pageNumber, searchTerms, enabledOnly)
}

func (s *service) Get(ctx context.Context, id uuid.UUID) (*Integration, error) {
	return s.repository.FindByID(ctx, id)
}

func (s *service) Save(ctx context.Context, data *Integration) error {
	driver := s.findDriver(data)
	if err := newValidator(s.repository, driver).validate(ctx, data); err != nil {
		return err
	}

	return s.repository.Save(ctx, data)
}

func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	inUse, err := s.repository.InUseByID(ctx, id)
	if err != nil {
		return err
	}

	if *inUse {
		return coreerror.New(i18n.M(ctx, i18n.K.CoreIntegrationInUse), true)
	}

	return s.repository.DeleteByID(ctx, id)
}

func (s *service) Exists(ctx context.Context, id uuid.UUID) (*bool, error) {
	return s.repository.ExistsByID(ctx, id)
}

func (s *service) ListOptions(
	ctx context.Context,
	integrationID uuid.UUID,
	pageNumber, pageSize int,
	searchTerms *string,
	tcpOnly bool,
) (*pagination.Page[DriverOption], error) {
	data, err := s.repository.FindByID(ctx, integrationID)
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, coreerror.New(
			i18n.M(ctx, i18n.K.CoreIntegrationNotFound),
			true,
		)
	}

	if !data.Enabled {
		return nil, coreerror.New(
			i18n.M(ctx, i18n.K.CoreIntegrationDisabled),
			true,
		)
	}

	driver := s.findDriver(data)
	if driver == nil {
		return nil, coreerror.New(
			i18n.M(ctx, i18n.K.CoreIntegrationNotFound),
			true,
		)
	}

	options, err := driver.GetAvailableOptions(
		ctx,
		data.Parameters,
		pageNumber,
		pageSize,
		searchTerms,
		tcpOnly,
	)
	if err != nil {
		return nil, err
	}

	sort.Slice(options.Contents, func(left, right int) bool {
		return options.Contents[left].Name < options.Contents[right].Name
	})

	return options, nil
}

func (s *service) GetOption(
	ctx context.Context,
	integrationID uuid.UUID,
	optionID string,
) (*DriverOption, error) {
	data, err := s.repository.FindByID(ctx, integrationID)
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, coreerror.New(
			i18n.M(ctx, i18n.K.CoreIntegrationNotFound),
			true,
		)
	}

	driver := s.findDriver(data)
	if driver == nil {
		return nil, coreerror.New(
			i18n.M(ctx, i18n.K.CoreIntegrationNotFound),
			true,
		)
	}

	if !data.Enabled {
		return nil, coreerror.New(
			i18n.M(ctx, i18n.K.CoreIntegrationDisabled),
			true,
		)
	}

	return driver.GetAvailableOptionByID(ctx, data.Parameters, optionID)
}

func (s *service) GetOptionURL(
	ctx context.Context,
	integrationID uuid.UUID,
	optionID string,
) (*string, []string, error) {
	data, err := s.repository.FindByID(ctx, integrationID)
	if err != nil {
		return nil, nil, err
	}

	if data == nil {
		return nil, nil, coreerror.New(
			i18n.M(ctx, i18n.K.CoreIntegrationNotFound),
			true,
		)
	}

	driver := s.findDriver(data)
	if driver == nil {
		return nil, nil, coreerror.New(
			i18n.M(ctx, i18n.K.CoreIntegrationNotFound),
			true,
		)
	}

	if !data.Enabled {
		return nil, nil, coreerror.New(
			i18n.M(ctx, i18n.K.CoreIntegrationDisabled),
			true,
		)
	}

	return driver.GetOptionProxyURL(ctx, data.Parameters, optionID)
}

func (s *service) GetAvailableDrivers(ctx context.Context) ([]AvailableDriver, error) {
	drivers := s.drivers()
	sort.Slice(drivers, func(left, right int) bool {
		return drivers[left].Name(ctx).String() < drivers[right].Name(ctx).String()
	})

	output := make([]AvailableDriver, len(drivers))
	for index, driver := range drivers {
		output[index] = AvailableDriver{
			ID:                  driver.ID(),
			Name:                driver.Name(ctx),
			Description:         driver.Description(ctx),
			ConfigurationFields: driver.ConfigurationFields(ctx),
		}
	}

	return output, nil
}

func (s *service) findDriver(data *Integration) Driver {
	for _, driver := range s.drivers() {
		if driver.ID() == data.Driver {
			return driver
		}
	}

	return nil
}
