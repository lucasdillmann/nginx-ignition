package server

import (
	"context"

	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/common/broadcast"
	"dillmann.com.br/nginx-ignition/core/common/coreerror"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/log"
	"dillmann.com.br/nginx-ignition/core/common/pagination"
)

type service struct {
	repository Repository
	providers  func() []Provider
}

func newService(
	repository Repository,
	providers func() []Provider,
) *service {
	return &service{
		repository: repository,
		providers:  providers,
	}
}

func (s *service) deleteById(ctx context.Context, id uuid.UUID) error {
	inUse, err := s.repository.IsInUseByID(ctx, id)
	if err != nil {
		return err
	}

	if inUse {
		return coreerror.New(
			"Certificate is being used by at least one host and/or global binding. Please update them and try again.",
			true,
		)
	}

	return s.repository.DeleteByID(ctx, id)
}

func (s *service) getById(ctx context.Context, id uuid.UUID) (*Certificate, error) {
	cert, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if cert != nil {
		availableProviders, err := s.availableProviders(ctx)
		if err != nil {
			return nil, err
		}

		provider := providerById(availableProviders, cert.ProviderID)
		dynamicfields.RemoveSensitiveFields(&cert.Parameters, provider.DynamicFields())
	}

	return cert, nil
}

func (s *service) list(ctx context.Context, pageSize, pageNumber int, searchTerms *string) (*pagination.Page[*Certificate], error) {
	certs, err := s.repository.FindPage(ctx, pageSize, pageNumber, searchTerms)
	if err != nil {
		return nil, err
	}

	availableProviders, err := s.availableProviders(ctx)
	if err != nil {
		return nil, err
	}

	for _, cert := range certs.Contents {
		provider := providerById(availableProviders, cert.ProviderID)
		dynamicfields.RemoveSensitiveFields(&cert.Parameters, provider.DynamicFields())
	}

	return certs, nil
}

func (s *service) availableProviders(_ context.Context) ([]*AvailableProvider, error) {
	var availableProviders []*AvailableProvider
	for _, provider := range s.providers() {
		availableProviders = append(availableProviders, &AvailableProvider{
			provider: provider,
		})
	}

	return availableProviders, nil
}

func (s *service) renewAllDue(ctx context.Context) error {
	certificates, err := s.repository.FindAllDueToRenew(ctx)
	if err != nil {
		return err
	}

	if len(certificates) == 0 {
		log.Infof("ServerCertificates auto-renew triggered, but no certificates are due to be renewed yet")
		return nil
	}

	for _, certificate := range certificates {
		err = s.renew(ctx, certificate.ID)
		if err != nil {
			log.Warnf("Error renewing certificate %s: %s", certificate.ID, err)
			continue
		}

		log.Infof("Certificate %s renewed successfully", certificate.ID)
	}

	broadcast.SendSignal(ctx, "core:nginx:reload")

	log.Infof("ServerCertificates auto-renew complemeted, %d were renewed", len(certificates))
	return nil
}

func (s *service) renew(ctx context.Context, id uuid.UUID) error {
	certificate, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return err
	}

	providers, err := s.availableProviders(ctx)
	if err != nil {
		return err
	}

	provider := providerById(providers, certificate.ProviderID)
	if provider == nil {
		return coreerror.New("Provider not found", true)
	}

	certificate, err = provider.Renew(ctx, certificate)
	if err != nil {
		return err
	}

	certificate.ID = id
	return s.repository.Save(ctx, certificate)
}

func (s *service) issue(ctx context.Context, request *IssueRequest) (*Certificate, error) {
	providers, err := s.availableProviders(ctx)
	if err != nil {
		return nil, err
	}

	provider := providerById(providers, request.ProviderID)
	if provider == nil {
		return nil, coreerror.New("Provider not found", true)
	}

	certificate, err := provider.Issue(ctx, request)
	if err != nil {
		return nil, err
	}

	err = s.repository.Save(ctx, certificate)
	if err != nil {
		return nil, err
	}

	return certificate, nil
}

func providerById(availableProviders []*AvailableProvider, id string) Provider {
	for _, p := range availableProviders {
		if p.provider.ID() == id {
			return p.provider
		}
	}

	return nil
}
