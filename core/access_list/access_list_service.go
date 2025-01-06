package access_list

import (
	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"dillmann.com.br/nginx-ignition/core/host"
	"errors"
	"github.com/google/uuid"
)

type service struct {
	repository     *Repository
	hostRepository *host.Repository
}

func (s *service) save(accessList *AccessList) error {
	if err := newValidator().validate(accessList); err != nil {
		return err
	}

	return (*s.repository).Save(accessList)
}

func (s *service) deleteById(id uuid.UUID) error {
	inUse, err := (*s.hostRepository).ExistsByAccessListId(id)
	if err != nil {
		return err
	}

	if inUse {
		return errors.New("access List is in use by one or more hosts")
	}

	return (*s.repository).DeleteById(id)
}

func (s *service) findById(id uuid.UUID) (*AccessList, error) {
	return (*s.repository).FindById(id)
}

func (s *service) findAll() (*[]AccessList, error) {
	return (*s.repository).FindAll()
}

func (s *service) list(
	pageNumber int64,
	pageSize int64,
	searchTerms string,
) (*pagination.Page[AccessList], error) {
	return (*s.repository).FindPage(pageNumber, pageSize, searchTerms)
}
