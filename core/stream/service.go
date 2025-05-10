package stream

import (
	"context"
	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"github.com/google/uuid"
)

type service struct {
	streamRepository Repository
}

func newService(streamRepository Repository) *service {
	return &service{streamRepository}
}

func (s *service) save(ctx context.Context, input *Stream) error {
	if err := newValidator().validate(ctx, input); err != nil {
		return err
	}

	return s.streamRepository.Save(ctx, input)
}

func (s *service) deleteByID(ctx context.Context, id uuid.UUID) error {
	return s.streamRepository.DeleteByID(ctx, id)
}

func (s *service) list(ctx context.Context, pageSize, pageNumber int, searchTerms *string) (*pagination.Page[*Stream], error) {
	return s.streamRepository.FindPage(ctx, pageSize, pageNumber, searchTerms)
}

func (s *service) getByID(ctx context.Context, id uuid.UUID) (*Stream, error) {
	return s.streamRepository.FindByID(ctx, id)
}

func (s *service) getAllEnabled(ctx context.Context) ([]*Stream, error) {
	return s.streamRepository.FindAllEnabled(ctx)
}

func (s *service) existsByID(ctx context.Context, id uuid.UUID) (bool, error) {
	return s.streamRepository.ExistsByID(ctx, id)
}
