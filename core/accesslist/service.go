package accesslist

import (
	"context"

	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/common/coreerror"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
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

func (s *service) Save(ctx context.Context, accessList *AccessList) error {
	if err := newValidator().validate(ctx, accessList); err != nil {
		return err
	}

	return s.repository.Save(ctx, accessList)
}

func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	inUse, err := s.repository.InUseByID(ctx, id)
	if err != nil {
		return err
	}

	if inUse {
		return coreerror.New(i18n.M(ctx, i18n.K.CoreAccesslistInUse), true)
	}

	return s.repository.DeleteByID(ctx, id)
}

func (s *service) Get(ctx context.Context, id uuid.UUID) (*AccessList, error) {
	return s.repository.FindByID(ctx, id)
}

func (s *service) GetAll(ctx context.Context) ([]AccessList, error) {
	return s.repository.FindAll(ctx)
}

func (s *service) List(
	ctx context.Context,
	pageSize,
	pageNumber int,
	searchTerms *string,
) (*pagination.Page[AccessList], error) {
	return s.repository.FindPage(ctx, pageNumber, pageSize, searchTerms)
}

func (s *service) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	return s.repository.ExistsByID(ctx, id)
}
