package access_list

import (
	"context"

	"dillmann.com.br/nginx-ignition/core/common/core_error"
	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"dillmann.com.br/nginx-ignition/core/host"
	"github.com/google/uuid"
)

type service struct {
	accessListRepository Repository
	hostRepository       host.Repository
}

func newService(accessListRepository Repository, hostRepository host.Repository) *service {
	return &service{
		accessListRepository: accessListRepository,
		hostRepository:       hostRepository,
	}
}

func (s *service) save(ctx context.Context, accessList *AccessList) error {
	if err := newValidator().validate(accessList); err != nil {
		return err
	}

	return s.accessListRepository.Save(ctx, accessList)
}

func (s *service) deleteById(ctx context.Context, id uuid.UUID) error {
	inUse, err := s.hostRepository.ExistsByAccessListID(ctx, id)
	if err != nil {
		return err
	}

	if inUse {
		return core_error.New("Access list is in use by one or more hosts", true)
	}

	return s.accessListRepository.DeleteByID(ctx, id)
}

func (s *service) findById(ctx context.Context, id uuid.UUID) (*AccessList, error) {
	return s.accessListRepository.FindByID(ctx, id)
}

func (s *service) findAll(ctx context.Context) ([]*AccessList, error) {
	return s.accessListRepository.FindAll(ctx)
}

func (s *service) list(
	ctx context.Context,
	pageSize,
	pageNumber int,
	searchTerms *string,
) (*pagination.Page[*AccessList], error) {
	return s.accessListRepository.FindPage(ctx, pageNumber, pageSize, searchTerms)
}
