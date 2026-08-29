package host

import (
	"context"

	"github.com/google/uuid"

	"nginx-ignition/internal/core/accesslist"
	"nginx-ignition/internal/core/binding"
	"nginx-ignition/internal/core/cache"
	"nginx-ignition/internal/core/certificate"
	"nginx-ignition/internal/core/common/pagination"
	"nginx-ignition/internal/core/integration"
	"nginx-ignition/internal/core/vpn"
)

type service struct {
	repository          Repository
	integrationCommands integration.Commands
	vpnCommands         vpn.Commands
	accessListCommands  accesslist.Commands
	cacheCommands       cache.Commands
	bindingCommands     binding.Commands
	certificateCommands certificate.Commands
}

func newCommands(
	repository Repository,
	integrationCommands integration.Commands,
	vpnCommands vpn.Commands,
	accessListCommands accesslist.Commands,
	cacheCommands cache.Commands,
	bindingCommands binding.Commands,
	certificateCommands certificate.Commands,
) Commands {
	return &service{
		repository:          repository,
		integrationCommands: integrationCommands,
		vpnCommands:         vpnCommands,
		accessListCommands:  accessListCommands,
		cacheCommands:       cacheCommands,
		bindingCommands:     bindingCommands,
		certificateCommands: certificateCommands,
	}
}

func (s *service) Save(ctx context.Context, input *Host) error {
	validatorInstance := newValidator(
		s.repository,
		s.integrationCommands,
		s.vpnCommands,
		s.accessListCommands,
		s.cacheCommands,
		s.bindingCommands,
		s.certificateCommands,
	)

	if err := validatorInstance.validate(ctx, input); err != nil {
		return err
	}

	return s.repository.Save(ctx, input)
}

func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repository.DeleteByID(ctx, id)
}

func (s *service) List(
	ctx context.Context,
	pageSize, pageNumber int,
	searchTerms *string,
) (*pagination.Page[Host], error) {
	return s.repository.FindPage(ctx, pageSize, pageNumber, searchTerms)
}

func (s *service) Get(ctx context.Context, id uuid.UUID) (*Host, error) {
	return s.repository.FindByID(ctx, id)
}

func (s *service) GetAllEnabled(ctx context.Context) ([]Host, error) {
	return s.repository.FindAllEnabled(ctx)
}

func (s *service) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	return s.repository.ExistsByID(ctx, id)
}
