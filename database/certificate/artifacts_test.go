package certificate

import (
	"time"

	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/certificate"
)

func newCertificate() *certificate.Certificate {
	return &certificate.Certificate{
		ID:         uuid.New(),
		IssuedAt:   time.Now(),
		ValidUntil: time.Now().Add(24 * time.Hour),
		ValidFrom:  time.Now(),
		RenewAfter: new(time.Now().Add(12 * time.Hour)),
		Parameters: map[string]any{
			"foo": "bar",
		},
		Metadata:   new("Test Metadata"),
		ProviderID: "test-provider",
		PrivateKey: "-----BEGIN PRIVATE KEY-----\n...\n-----END PRIVATE KEY-----",
		PublicKey:  "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
		DomainNames: []string{
			"example.com",
			"sub.example.com",
		},
		CertificationChain: []string{
			"-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
		},
	}
}
