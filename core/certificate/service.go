package certificate

import (
	"dillmann.com.br/nginx-ignition/core/common/core_error"
	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"dillmann.com.br/nginx-ignition/core/host"
	"dillmann.com.br/nginx-ignition/core/settings"
	"github.com/google/uuid"
	"go.uber.org/dig"
)

type service struct {
	container             *dig.Container
	certificateRepository *Repository
	hostRepository        *host.Repository
	settingsRepository    *settings.Repository
	providers             *[]*Provider
}

func newService(
	container *dig.Container,
	certificateRepository *Repository,
	hostRepository *host.Repository,
	settingsRepository *settings.Repository,
) *service {
	return &service{
		container:             container,
		certificateRepository: certificateRepository,
		hostRepository:        hostRepository,
		settingsRepository:    settingsRepository,
		providers:             nil,
	}
}

func (s *service) deleteById(id uuid.UUID) error {
	inUse, err := (*s.hostRepository).ExistsByCertificateID(id)
	if err != nil {
		return err
	}

	if inUse {
		return core_error.New(
			"Certificate is being used by at least one host. Please update them and try again.",
			true,
		)
	}

	cfg, err := (*s.settingsRepository).Get()
	if err != nil {
		return err
	}

	for _, binding := range cfg.GlobalBindings {
		if binding.CertificateID != nil && *binding.CertificateID == id {
			inUse = true
			break
		}
	}

	if inUse {
		return core_error.New(
			"Certificate is being used by a global binding. Please update the settings and try again.",
			true,
		)
	}

	return (*s.certificateRepository).DeleteByID(id)
}

func (s *service) getById(id uuid.UUID) (*Certificate, error) {
	return (*s.certificateRepository).FindByID(id)
}

func (s *service) list(pageSize, pageNumber int, searchTerms *string) (*pagination.Page[*Certificate], error) {
	return (*s.certificateRepository).FindPage(pageSize, pageNumber, searchTerms)
}

func (s *service) availableProviders() ([]*AvailableProvider, error) {
	return []*AvailableProvider{}, nil
}

func (s *service) renew(_ uuid.UUID) error {
	return core_error.New("Not implemented yet", false)
}

func (s *service) issue(_ *IssueRequest) error {
	return core_error.New("Not implemented yet", false)
}
