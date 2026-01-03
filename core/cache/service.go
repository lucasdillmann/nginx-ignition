package cache

import (
	"context"

	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/common/coreerror"
	"dillmann.com.br/nginx-ignition/core/common/pagination"
)

type service struct {
	repository Repository
}

func newCommands(repository Repository) Commands {
	return &service{
		repository: repository,
	}
}

func (s *service) Save(ctx context.Context, c *Cache) error {
	if err := newValidator().validate(c); err != nil {
		return err
	}

	return s.repository.Save(ctx, c)
}

func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	inUse, err := s.repository.InUseByID(ctx, id)
	if err != nil {
		return err
	}

	if inUse {
		return coreerror.New("Cache configuration is in use by one or more hosts", true)
	}

	return s.repository.DeleteByID(ctx, id)
}

func (s *service) Get(ctx context.Context, id uuid.UUID) (*Cache, error) {
	return s.repository.FindByID(ctx, id)
}

func (s *service) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	return s.repository.ExistsByID(ctx, id)
}

func (s *service) List(
	ctx context.Context,
	pageSize,
	pageNumber int,
	searchTerms *string,
) (*pagination.Page[Cache], error) {
	return s.repository.FindPage(ctx, pageNumber, pageSize, searchTerms)
}

func (s *service) GetAllInUse(ctx context.Context) ([]Cache, error) {
	return s.repository.FindAllInUse(ctx)
}
