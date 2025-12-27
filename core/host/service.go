package host

import (
	"context"

	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/accesslist"
	"dillmann.com.br/nginx-ignition/core/binding"
	"dillmann.com.br/nginx-ignition/core/cache"
	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"dillmann.com.br/nginx-ignition/core/integration"
	"dillmann.com.br/nginx-ignition/core/vpn"
)

type service struct {
	repository          Repository
	integrationCommands *integration.Commands
	vpnCommands         *vpn.Commands
	accessListCommands  *accesslist.Commands
	cacheCommands       *cache.Commands
	bindingCommands     *binding.Commands
}

func newService(
	repository Repository,
	integrationCommands *integration.Commands,
	vpnCommands *vpn.Commands,
	accessListCommands *accesslist.Commands,
	cacheCommands *cache.Commands,
	bindingCommands *binding.Commands,
) *service {
	return &service{
		repository:          repository,
		integrationCommands: integrationCommands,
		vpnCommands:         vpnCommands,
		accessListCommands:  accessListCommands,
		cacheCommands:       cacheCommands,
		bindingCommands:     bindingCommands,
	}
}

func (s *service) save(ctx context.Context, input *Host) error {
	validatorInstance := newValidator(
		s.repository,
		s.integrationCommands,
		s.vpnCommands,
		s.accessListCommands,
		s.cacheCommands,
		s.bindingCommands,
	)

	if err := validatorInstance.validate(ctx, input); err != nil {
		return err
	}

	return s.repository.Save(ctx, input)
}

func (s *service) deleteByID(ctx context.Context, id uuid.UUID) error {
	return s.repository.DeleteByID(ctx, id)
}

func (s *service) list(ctx context.Context, pageSize, pageNumber int, searchTerms *string) (*pagination.Page[Host], error) {
	return s.repository.FindPage(ctx, pageSize, pageNumber, searchTerms)
}

func (s *service) getByID(ctx context.Context, id uuid.UUID) (*Host, error) {
	return s.repository.FindByID(ctx, id)
}

func (s *service) getAllEnabled(ctx context.Context) ([]Host, error) {
	return s.repository.FindAllEnabled(ctx)
}

func (s *service) existsByID(ctx context.Context, id uuid.UUID) (bool, error) {
	return s.repository.ExistsByID(ctx, id)
}
