package selfsigned

import (
	"context"
	"crypto/x509"
	"encoding/base64"
	"time"

	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/certificate/commons"
	"dillmann.com.br/nginx-ignition/core/certificate"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

type Provider struct{}

func New() *Provider {
	return &Provider{}
}

func (p *Provider) ID() string {
	return "SELF_SIGNED"
}

func (p *Provider) Name() string {
	return "Self-signed certificate"
}

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return []*dynamic_fields.DynamicField{}
}

func (p *Provider) Priority() int {
	return 3
}

func (p *Provider) Issue(_ context.Context, request *certificate.IssueRequest) (*certificate.Certificate, error) {
	if err := commons.Validate(request, validationRules{}); err != nil {
		return nil, err
	}

	certPEM, keyPEM, err := buildPEMs(request.DomainNames)
	if err != nil {
		return nil, err
	}

	renewAfter := time.Now().Add(335 * 24 * time.Hour)
	cert := &certificate.Certificate{
		ID:                 uuid.New(),
		DomainNames:        request.DomainNames,
		ProviderID:         p.ID(),
		IssuedAt:           time.Now(),
		ValidFrom:          time.Now().Add(-1 * time.Hour),
		ValidUntil:         time.Now().Add(365 * 24 * time.Hour),
		RenewAfter:         &renewAfter,
		PrivateKey:         *keyPEM,
		PublicKey:          *certPEM,
		CertificationChain: []string{},
		Parameters:         map[string]any{},
		Metadata:           nil,
	}

	return cert, nil
}

func (p *Provider) Renew(_ context.Context, current *certificate.Certificate) (*certificate.Certificate, error) {
	certPEM, keyPEM, err := buildPEMs(current.DomainNames)
	if err != nil {
		return nil, err
	}

	renewAfter := time.Now().Add(335 * 24 * time.Hour)
	cert := &certificate.Certificate{
		ID:                 current.ID,
		DomainNames:        current.DomainNames,
		ProviderID:         p.ID(),
		IssuedAt:           time.Now(),
		ValidFrom:          time.Now().Add(-1 * time.Hour),
		ValidUntil:         time.Now().Add(365 * 24 * time.Hour),
		RenewAfter:         &renewAfter,
		PrivateKey:         *keyPEM,
		PublicKey:          *certPEM,
		CertificationChain: []string{},
		Parameters:         current.Parameters,
		Metadata:           current.Metadata,
	}

	return cert, nil
}

func buildPEMs(domainNames []string) (*string, *string, error) {
	factory, err := newFactory()
	if err != nil {
		return nil, nil, err
	}

	cert, key, err := factory.build(domainNames)
	if err != nil {
		return nil, nil, err
	}

	certBase64 := base64.StdEncoding.EncodeToString(cert.Raw)
	keyBase64 := base64.StdEncoding.EncodeToString(x509.MarshalPKCS1PrivateKey(key))

	return &certBase64, &keyBase64, nil
}
