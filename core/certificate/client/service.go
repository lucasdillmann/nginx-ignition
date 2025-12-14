package client

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
	return &service{repository}
}

func (s *service) list(
	ctx context.Context,
	pageSize, pageNumber int,
	searchTerms *string,
) (*pagination.Page[*Certificate], error) {
	return s.repository.FindPage(ctx, pageSize, pageNumber, searchTerms)
}

func (s *service) getById(ctx context.Context, id uuid.UUID) (*Certificate, error) {
	return s.repository.FindByID(ctx, id)
}

func (s *service) deleteById(ctx context.Context, id uuid.UUID) error {
	inUse, err := s.repository.IsInUseByID(ctx, id)
	if err != nil {
		return err
	}

	if inUse {
		return coreerror.New(
			"Client certificate is being used by at least one host. Please update them and try again.",
			true,
		)
	}

	return s.repository.DeleteByID(ctx, id)
}

// TODO: Implement remaining functions (as defined in the Commands struct)
