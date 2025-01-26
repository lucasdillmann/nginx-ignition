package selfsigned

import (
	"crypto/x509"
	"dillmann.com.br/nginx-ignition/core/certificate"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
	"encoding/base64"
	"github.com/google/uuid"
	"time"
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

func (p *Provider) Issue(request *certificate.IssueRequest) (*certificate.Certificate, error) {
	dereferencedNames := dereference(request.DomainNames)
	certPEM, keyPEM, err := buildPEMs(dereferencedNames)
	if err != nil {
		return nil, err
	}

	renewAfter := time.Now().Add(335 * 24 * time.Hour)
	cert := &certificate.Certificate{
		ID:                 uuid.New(),
		DomainNames:        dereferencedNames,
		ProviderID:         p.ID(),
		IssuedAt:           time.Now(),
		ValidFrom:          time.Now().Add(-1 * time.Hour),
		ValidUntil:         time.Now().Add(365 * 24 * time.Hour),
		RenewAfter:         &renewAfter,
		PrivateKey:         *keyPEM,
		PublicKey:          *certPEM,
		CertificationChain: []string{},
		Parameters:         map[string]interface{}{},
		Metadata:           nil,
	}

	return cert, nil
}

func (p *Provider) Renew(current *certificate.Certificate) (*certificate.Certificate, error) {
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

func dereference(domainNames []*string) []string {
	dereferencedNames := make([]string, len(domainNames))
	for i, name := range domainNames {
		dereferencedNames[i] = *name
	}
	return dereferencedNames
}
