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
	var providers *[]Provider

	if err := s.container.Invoke(func(p *[]Provider) {
		providers = p
	}); err != nil {
		return nil, err
	}

	if providers == nil {
		return nil, core_error.New("No providers found", false)
	}

	var availableProviders []*AvailableProvider
	for _, provider := range *providers {
		availableProviders = append(availableProviders, &AvailableProvider{
			provider: &provider,
		})
	}

	return availableProviders, nil
}

func (s *service) renew(certificateId uuid.UUID) error {
	certificate, err := (*s.certificateRepository).FindByID(certificateId)
	if err != nil {
		return err
	}

	providers, err := s.availableProviders()
	if err != nil {
		return err
	}

	provider := providerById(providers, certificate.ProviderID)
	if provider == nil {
		return core_error.New("Provider not found", true)
	}

	certificate, err = (*provider).Renew(certificate)
	if err != nil {
		return err
	}

	certificate.ID = certificateId
	return (*s.certificateRepository).Save(certificate)
}

func (s *service) issue(request *IssueRequest) (*Certificate, error) {
	providers, err := s.availableProviders()
	if err != nil {
		return nil, err
	}

	provider := providerById(providers, request.ProviderID)
	if provider == nil {
		return nil, core_error.New("Provider not found", true)
	}

	certificate, err := (*provider).Issue(request)
	if err != nil {
		return nil, err
	}

	err = (*s.certificateRepository).Save(certificate)
	if err != nil {
		return nil, err
	}

	return certificate, nil
}

func providerById(availableProviders []*AvailableProvider, id string) *Provider {
	for _, p := range availableProviders {
		if (*p.provider).ID() == id {
			return p.provider
		}
	}

	return nil
}
