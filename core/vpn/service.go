package vpn

import (
	"context"
	"sort"

	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/common/core_error"
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
) (*pagination.Page[*VPN], error) {
	return s.repository.FindPage(ctx, pageSize, pageNumber, searchTerms, enabledOnly)
}

func (s *service) getById(ctx context.Context, id uuid.UUID) (*VPN, error) {
	return s.repository.FindByID(ctx, id)
}

func (s *service) save(ctx context.Context, data *VPN) error {
	driver := s.findDriver(data)
	if err := newValidator(s.repository, driver).validate(ctx, data); err != nil {
		return err
	}

	return s.repository.Save(ctx, data)
}

func (s *service) deleteById(ctx context.Context, id uuid.UUID) error {
	inUse, err := s.repository.InUseByID(ctx, id)
	if err != nil {
		return err
	}

	if *inUse {
		return core_error.New("VPN is in use by one or more hosts", true)
	}

	return s.repository.DeleteByID(ctx, id)
}

func (s *service) existsById(ctx context.Context, id uuid.UUID) (*bool, error) {
	return s.repository.ExistsByID(ctx, id)
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
			ConfigurationFields: driver.ConfigurationFields(),
		}
	}

	return &output, nil
}

func (s *service) findDriver(data *VPN) Driver {
	for _, driver := range s.drivers() {
		if driver.ID() == data.Driver {
			return driver
		}
	}

	return nil
}
