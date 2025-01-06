package access_list

import (
	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"dillmann.com.br/nginx-ignition/core/host"
	"errors"
	"github.com/google/uuid"
)

type accessListService struct {
	accessListRepository *AccessListRepository
	hostRepository       *host.HostRepository
}

func (service *accessListService) save(accessList *AccessList) error {
	err := newValidator().validate(accessList)
	if err != nil {
		return err
	}

	return (*service.accessListRepository).Save(accessList)
}

func (service *accessListService) deleteById(id uuid.UUID) error {
	inUse, err := (*service.hostRepository).ExistsByAccessListId(id)
	if err != nil {
		return err
	}

	if inUse {
		return errors.New("access List is in use by one or more hosts")
	}

	return (*service.accessListRepository).DeleteById(id)
}

func (service *accessListService) findById(id uuid.UUID) (*AccessList, error) {
	return (*service.accessListRepository).FindById(id)
}

func (service *accessListService) findAll() (*[]AccessList, error) {
	return (*service.accessListRepository).FindAll()
}

func (service *accessListService) list(
	pageNumber int64,
	pageSize int64,
	searchTerms string,
) (*pagination.Page[AccessList], error) {
	return (*service.accessListRepository).FindPage(pageNumber, pageSize, searchTerms)
}
