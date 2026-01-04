package certificate

import (
	"time"

	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/certificate"
	"dillmann.com.br/nginx-ignition/core/common/pagination"
)

func newCertificate() *certificate.Certificate {
	return &certificate.Certificate{
		ID:          uuid.New(),
		ProviderID:  "custom",
		DomainNames: []string{"example.com"},
		IssuedAt:    time.Now(),
		ValidFrom:   time.Now(),
		ValidUntil:  time.Now().AddDate(0, 3, 0),
		PrivateKey:  "private-key",
		PublicKey:   "public-key",
	}
}

func newCertificatePage() *pagination.Page[certificate.Certificate] {
	return pagination.Of([]certificate.Certificate{
		*newCertificate(),
	})
}

func newIssueCertificateRequest() *issueCertificateRequest {
	return &issueCertificateRequest{
		ProviderID:  "custom",
		DomainNames: []string{"example.com"},
		Parameters: map[string]any{
			"param1": "value1",
		},
	}
}
