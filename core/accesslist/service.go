package accesslist

import (
	"context"

	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/common/coreerror"
	"dillmann.com.br/nginx-ignition/core/common/pagination"
)

type service struct {
	repository Repository
}

func newService(repository Repository) *service {
	return &service{
		repository: repository,
	}
}

func (s *service) save(ctx context.Context, accessList *AccessList) error {
	if err := newValidator().validate(accessList); err != nil {
		return err
	}

	return s.repository.Save(ctx, accessList)
}

func (s *service) deleteById(ctx context.Context, id uuid.UUID) error {
	inUse, err := s.repository.InUseByID(ctx, id)
	if err != nil {
		return err
	}

	if inUse {
		return coreerror.New("Access list is in use by one or more hosts", true)
	}

	return s.repository.DeleteByID(ctx, id)
}

func (s *service) findById(ctx context.Context, id uuid.UUID) (*AccessList, error) {
	return s.repository.FindByID(ctx, id)
}

func (s *service) findAll(ctx context.Context) ([]AccessList, error) {
	return s.repository.FindAll(ctx)
}

func (s *service) list(
	ctx context.Context,
	pageSize,
	pageNumber int,
	searchTerms *string,
) (*pagination.Page[AccessList], error) {
	return s.repository.FindPage(ctx, pageNumber, pageSize, searchTerms)
}

func (s *service) existsByID(ctx context.Context, id uuid.UUID) (bool, error) {
	return s.repository.ExistsByID(ctx, id)
}
