package host

import (
	"context"

	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/accesslist"
	"dillmann.com.br/nginx-ignition/core/cache"
	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"dillmann.com.br/nginx-ignition/core/common/validation"
	"dillmann.com.br/nginx-ignition/core/integration"
	"dillmann.com.br/nginx-ignition/core/vpn"
)

type service struct {
	hostRepository      Repository
	integrationCommands *integration.Commands
	vpnCommands         *vpn.Commands
	accessListCommands  *accesslist.Commands
	cacheCommands       *cache.Commands
}

func newService(
	hostRepository Repository,
	integrationCommands *integration.Commands,
	vpnCommands *vpn.Commands,
	accessListCommands *accesslist.Commands,
	cacheCommands *cache.Commands,
) *service {
	return &service{
		hostRepository,
		integrationCommands,
		vpnCommands,
		accessListCommands,
		cacheCommands,
	}
}

func (s *service) save(ctx context.Context, input *Host) error {
	validatorInstance := newValidator(
		s.hostRepository,
		s.integrationCommands,
		s.vpnCommands,
		s.accessListCommands,
		s.cacheCommands,
	)

	if err := validatorInstance.validate(ctx, input); err != nil {
		return err
	}

	return s.hostRepository.Save(ctx, input)
}

func (s *service) deleteByID(ctx context.Context, id uuid.UUID) error {
	return s.hostRepository.DeleteByID(ctx, id)
}

func (s *service) list(ctx context.Context, pageSize, pageNumber int, searchTerms *string) (*pagination.Page[*Host], error) {
	return s.hostRepository.FindPage(ctx, pageSize, pageNumber, searchTerms)
}

func (s *service) getByID(ctx context.Context, id uuid.UUID) (*Host, error) {
	return s.hostRepository.FindByID(ctx, id)
}

func (s *service) getAllEnabled(ctx context.Context) ([]*Host, error) {
	return s.hostRepository.FindAllEnabled(ctx)
}

func (s *service) existsByID(ctx context.Context, id uuid.UUID) (bool, error) {
	return s.hostRepository.ExistsByID(ctx, id)
}

func (s *service) validateBinding(
	ctx context.Context,
	path string,
	index int,
	binding *Binding,
	context *validation.ConsistencyValidator,
) error {
	validatorInstance := &validator{
		s.hostRepository,
		s.integrationCommands,
		s.vpnCommands,
		s.accessListCommands,
		s.cacheCommands,
		context,
	}
	return validatorInstance.validateBinding(ctx, path, binding, index)
}
