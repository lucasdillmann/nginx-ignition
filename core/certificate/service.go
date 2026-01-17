package certificate

import (
	"context"

	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/common/broadcast"
	"dillmann.com.br/nginx-ignition/core/common/coreerror"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
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

func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	inUse, err := s.repository.InUseByID(ctx, id)
	if err != nil {
		return err
	}

	if inUse {
		return coreerror.New(i18n.M(ctx, i18n.K.CertificateErrorInUse), true)
	}

	return s.repository.DeleteByID(ctx, id)
}

func (s *service) Get(ctx context.Context, id uuid.UUID) (*Certificate, error) {
	cert, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if cert != nil {
		availableProviders, err := s.AvailableProviders(ctx)
		if err != nil {
			return nil, err
		}

		provider := providerByID(availableProviders, cert.ProviderID)
		dynamicfields.RemoveSensitiveFields(&cert.Parameters, provider.DynamicFields(ctx))
	}

	return cert, nil
}

func (s *service) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	return s.repository.ExistsByID(ctx, id)
}

func (s *service) List(
	ctx context.Context,
	pageSize, pageNumber int,
	searchTerms *string,
) (*pagination.Page[Certificate], error) {
	certs, err := s.repository.FindPage(ctx, pageSize, pageNumber, searchTerms)
	if err != nil {
		return nil, err
	}

	availableProviders, err := s.AvailableProviders(ctx)
	if err != nil {
		return nil, err
	}

	for _, cert := range certs.Contents {
		provider := providerByID(availableProviders, cert.ProviderID)
		dynamicfields.RemoveSensitiveFields(&cert.Parameters, provider.DynamicFields(ctx))
	}

	return certs, nil
}

func (s *service) AvailableProviders(_ context.Context) ([]AvailableProvider, error) {
	availableProviders := make([]AvailableProvider, 0, len(s.providers()))
	for _, provider := range s.providers() {
		availableProviders = append(availableProviders, AvailableProvider{
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
		log.Infof(
			"Certificates auto-renew triggered, but no certificates are due to be renewed yet",
		)
		return nil
	}

	for _, certificate := range certificates {
		err = s.Renew(ctx, certificate.ID)
		if err != nil {
			log.Warnf("Error renewing certificate %s: %s", certificate.ID, err)
			continue
		}

		log.Infof("Certificate %s renewed successfully", certificate.ID)
	}

	broadcast.SendSignal(ctx, "core:nginx:reload")

	log.Infof("Certificates auto-renew complemeted, %d were renewed", len(certificates))
	return nil
}

func (s *service) Renew(ctx context.Context, certificateID uuid.UUID) error {
	certificate, err := s.repository.FindByID(ctx, certificateID)
	if err != nil {
		return err
	}

	providers, err := s.AvailableProviders(ctx)
	if err != nil {
		return err
	}

	provider := providerByID(providers, certificate.ProviderID)
	if provider == nil {
		return coreerror.New(i18n.M(ctx, i18n.K.CertificateErrorProviderNotFound), true)
	}

	certificate, err = provider.Renew(ctx, certificate)
	if err != nil {
		return err
	}

	certificate.ID = certificateID
	return s.repository.Save(ctx, certificate)
}

func (s *service) Issue(ctx context.Context, request *IssueRequest) (*Certificate, error) {
	providers, err := s.AvailableProviders(ctx)
	if err != nil {
		return nil, err
	}

	provider := providerByID(providers, request.ProviderID)
	if provider == nil {
		return nil, coreerror.New(i18n.M(ctx, i18n.K.CertificateErrorProviderNotFound), true)
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

func (s *service) autoRenewSettings(ctx context.Context) (*AutoRenewSettings, error) {
	return s.repository.GetAutoRenewSettings(ctx)
}

func providerByID(availableProviders []AvailableProvider, id string) Provider {
	for _, p := range availableProviders {
		if p.provider.ID() == id {
			return p.provider
		}
	}

	return nil
}
