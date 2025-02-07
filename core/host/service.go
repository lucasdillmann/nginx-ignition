package host

import (
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

func (s *service) save(input *Host) error {
	if err := newValidator(s.hostRepository).validate(input); err != nil {
		return err
	}

	return (*s.hostRepository).Save(input)
}

func (s *service) deleteByID(id uuid.UUID) error {
	return (*s.hostRepository).DeleteByID(id)
}

func (s *service) list(pageSize, pageNumber int, searchTerms *string) (*pagination.Page[*Host], error) {
	return (*s.hostRepository).FindPage(pageSize, pageNumber, searchTerms)
}

func (s *service) getByID(id uuid.UUID) (*Host, error) {
	return (*s.hostRepository).FindByID(id)
}

func (s *service) getAllEnabled() ([]*Host, error) {
	return (*s.hostRepository).FindAllEnabled()
}

func (s *service) existsByID(id uuid.UUID) (bool, error) {
	return (*s.hostRepository).ExistsByID(id)
}

func (s *service) validateBinding(
	path string,
	index int,
	binding *Binding,
	context *validation.ConsistencyValidator,
) error {
	validatorInstance := &validator{s.hostRepository, context}
	return validatorInstance.validateBinding(path, binding, index)
}
