package vpn

import (
	"context"
	"sort"

	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"dillmann.com.br/nginx-ignition/core/common/coreerror"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
	"dillmann.com.br/nginx-ignition/core/common/pagination"
)

type service struct {
	repository Repository
	cfg        *configuration.Configuration
	drivers    func() []Driver
}

func newService(
	cfg *configuration.Configuration,
	repository Repository,
	drivers func() []Driver,
) *service {
	return &service{
		cfg:        cfg,
		repository: repository,
		drivers:    drivers,
	}
}

func (s *service) List(
	ctx context.Context,
	pageSize, pageNumber int,
	searchTerms *string,
	enabledOnly bool,
) (*pagination.Page[VPN], error) {
	return s.repository.FindPage(ctx, pageSize, pageNumber, searchTerms, enabledOnly)
}

func (s *service) Get(ctx context.Context, id uuid.UUID) (*VPN, error) {
	return s.repository.FindByID(ctx, id)
}

func (s *service) Save(ctx context.Context, data *VPN) error {
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
		return coreerror.New(i18n.M(ctx, i18n.K.CoreVpnInUse), true)
	}

	return s.repository.DeleteByID(ctx, id)
}

func (s *service) Exists(ctx context.Context, id uuid.UUID) (*bool, error) {
	return s.repository.ExistsByID(ctx, id)
}

func (s *service) GetAvailableDrivers(ctx context.Context) ([]AvailableDriver, error) {
	drivers := s.drivers()
	sort.Slice(drivers, func(left, right int) bool {
		return drivers[left].Name(ctx).String() < drivers[right].Name(ctx).String()
	})

	output := make([]AvailableDriver, len(drivers))
	for index, driver := range drivers {
		output[index] = AvailableDriver{
			ID:                    driver.ID(),
			Name:                  driver.Name(ctx),
			ImportantInstructions: driver.ImportantInstructions(ctx),
			ConfigurationFields:   driver.ConfigurationFields(ctx),
		}
	}

	return output, nil
}

func (s *service) Start(ctx context.Context, endpoint Endpoint) error {
	data, driver, configDir, err := s.resolveValues(ctx, endpoint.VPNID())
	if err != nil {
		return err
	}

	return driver.Start(ctx, *configDir, endpoint, data.Parameters)
}

func (s *service) Reload(ctx context.Context, endpoint Endpoint) error {
	data, driver, configDir, err := s.resolveValues(ctx, endpoint.VPNID())
	if err != nil {
		return err
	}

	return driver.Reload(ctx, *configDir, endpoint, data.Parameters)
}

func (s *service) Stop(ctx context.Context, endpoint Endpoint) error {
	_, driver, _, err := s.resolveValues(ctx, endpoint.VPNID())
	if err != nil {
		return err
	}

	return driver.Stop(ctx, endpoint)
}

func (s *service) resolveValues(ctx context.Context, id uuid.UUID) (*VPN, Driver, *string, error) {
	data, err := s.Get(ctx, id)
	if err != nil {
		return nil, nil, nil, err
	}

	driver := s.findDriver(data)
	if driver == nil {
		return nil, nil, nil, coreerror.New(i18n.M(ctx, i18n.K.CoreVpnDriverNotFound), false)
	}

	configDir, err := s.cfg.Get("nginx-ignition.vpn.config-path")
	if err != nil {
		return nil, nil, nil, err
	}

	return data, driver, &configDir, nil
}

func (s *service) findDriver(data *VPN) Driver {
	for _, driver := range s.drivers() {
		if driver.ID() == data.Driver {
			return driver
		}
	}

	return nil
}
