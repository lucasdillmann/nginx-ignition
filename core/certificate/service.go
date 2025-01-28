package certificate

import (
	"dillmann.com.br/nginx-ignition/core/common/broadcast"
	"dillmann.com.br/nginx-ignition/core/common/core_error"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
	"dillmann.com.br/nginx-ignition/core/common/log"
	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"dillmann.com.br/nginx-ignition/core/host"
	"dillmann.com.br/nginx-ignition/core/settings"
	"github.com/google/uuid"
)

type service struct {
	certificateRepository *Repository
	hostRepository        *host.Repository
	settingsRepository    *settings.Repository
	providerResolver      func() ([]Provider, error)
}

func newService(
	certificateRepository *Repository,
	hostRepository *host.Repository,
	settingsRepository *settings.Repository,
	providerResolver func() ([]Provider, error),
) *service {
	return &service{
		certificateRepository: certificateRepository,
		hostRepository:        hostRepository,
		settingsRepository:    settingsRepository,
		providerResolver:      providerResolver,
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
	cert, err := (*s.certificateRepository).FindByID(id)
	if err != nil {
		return nil, err
	}

	if cert != nil {
		availableProviders, err := s.availableProviders()
		if err != nil {
			return nil, err
		}

		provider := providerById(availableProviders, cert.ProviderID)
		dynamic_fields.RemoveSensitiveFields(&cert.Parameters, (*provider).DynamicFields())
	}

	return cert, nil
}

func (s *service) list(pageSize, pageNumber int, searchTerms *string) (*pagination.Page[*Certificate], error) {
	certs, err := (*s.certificateRepository).FindPage(pageSize, pageNumber, searchTerms)
	if err != nil {
		return nil, err
	}

	availableProviders, err := s.availableProviders()
	if err != nil {
		return nil, err
	}

	for _, cert := range certs.Contents {
		provider := providerById(availableProviders, cert.ProviderID)
		dynamic_fields.RemoveSensitiveFields(&cert.Parameters, (*provider).DynamicFields())
	}

	return certs, nil
}

func (s *service) availableProviders() ([]*AvailableProvider, error) {
	providers, err := s.providerResolver()
	if err != nil {
		return nil, err
	}

	if providers == nil {
		return nil, core_error.New("No providers found", false)
	}

	var availableProviders []*AvailableProvider
	for _, provider := range providers {
		availableProviders = append(availableProviders, &AvailableProvider{
			provider: &provider,
		})
	}

	return availableProviders, nil
}

func (s *service) renewAllDue() error {
	certificates, err := (*s.certificateRepository).FindAllDueToRenew()
	if err != nil {
		return err
	}

	if len(certificates) == 0 {
		log.Infof("Certificates auto-renew triggered, but no certificates are due to be renewed yet")
		return nil
	}

	for _, certificate := range certificates {
		err = s.renew(certificate.ID)
		if err != nil {
			log.Warnf("Error renewing certificate %s: %s", certificate.ID, err)
			continue
		}

		log.Infof("Certificate %s renewed successfully", certificate.ID)
	}

	broadcast.SendSignal("core:nginx:reload")

	log.Infof("Certificates auto-renew complemeted, %d were renewed", len(certificates))
	return nil
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
