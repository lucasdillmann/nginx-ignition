package access_list

import (
	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"dillmann.com.br/nginx-ignition/core/host"
	"errors"
	"github.com/google/uuid"
)

type service struct {
	accessListRepository *Repository
	hostRepository       *host.Repository
}

func newService(accessListRepository *Repository, hostRepository *host.Repository) *service {
	return &service{
		accessListRepository: accessListRepository,
		hostRepository:       hostRepository,
	}
}

func (s *service) save(accessList *AccessList) error {
	if err := newValidator().validate(accessList); err != nil {
		return err
	}

	return (*s.accessListRepository).Save(accessList)
}

func (s *service) deleteById(id uuid.UUID) error {
	inUse, err := (*s.hostRepository).ExistsByAccessListID(id)
	if err != nil {
		return err
	}

	if inUse {
		return errors.New("access List is in use by one or more hosts")
	}

	return (*s.accessListRepository).DeleteByID(id)
}

func (s *service) findById(id uuid.UUID) (*AccessList, error) {
	return (*s.accessListRepository).FindByID(id)
}

func (s *service) findAll() ([]*AccessList, error) {
	return (*s.accessListRepository).FindAll()
}

func (s *service) list(
	pageSize,
	pageNumber int,
	searchTerms *string,
) (*pagination.Page[*AccessList], error) {
	return (*s.accessListRepository).FindPage(pageNumber, pageSize, searchTerms)
}
