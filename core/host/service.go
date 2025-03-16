package host

import (
	"context"
	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"dillmann.com.br/nginx-ignition/core/common/validation"
	"github.com/google/uuid"
)

type service struct {
	hostRepository *Repository
}

func newService(
	hostRepository *Repository,
) *service {
	return &service{
		hostRepository,
	}
}

func (s *service) save(ctx context.Context, input *Host) error {
	if err := newValidator(s.hostRepository).validate(ctx, input); err != nil {
		return err
	}

	return (*s.hostRepository).Save(ctx, input)
}

func (s *service) deleteByID(ctx context.Context, id uuid.UUID) error {
	return (*s.hostRepository).DeleteByID(ctx, id)
}

func (s *service) list(ctx context.Context, pageSize, pageNumber int, searchTerms *string) (*pagination.Page[*Host], error) {
	return (*s.hostRepository).FindPage(ctx, pageSize, pageNumber, searchTerms)
}

func (s *service) getByID(ctx context.Context, id uuid.UUID) (*Host, error) {
	return (*s.hostRepository).FindByID(ctx, id)
}

func (s *service) getAllEnabled(ctx context.Context) ([]*Host, error) {
	return (*s.hostRepository).FindAllEnabled(ctx)
}

func (s *service) existsByID(ctx context.Context, id uuid.UUID) (bool, error) {
	return (*s.hostRepository).ExistsByID(ctx, id)
}

func (s *service) validateBinding(
	ctx context.Context,
	path string,
	index int,
	binding *Binding,
	context *validation.ConsistencyValidator,
) error {
	validatorInstance := &validator{s.hostRepository, context}
	return validatorInstance.validateBinding(ctx, path, binding, index)
}
