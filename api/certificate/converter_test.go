package certificate

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"dillmann.com.br/nginx-ignition/core/certificate"
)

func Test_Converter(t *testing.T) {
	t.Run("toCertificateResponse", func(t *testing.T) {
		t.Run("converts certificate to response DTO", func(t *testing.T) {
			id := uuid.New()
			now := time.Now()
			input := &certificate.Certificate{
				ID:          id,
				DomainNames: []string{"example.com"},
				ProviderID:  "digitalocean",
				IssuedAt:    now,
				ValidUntil:  now,
				ValidFrom:   now,
				RenewAfter:  &now,
				Parameters:  map[string]any{"token": "secret"},
			}

			result := toCertificateResponse(input)

			assert.NotNil(t, result)
			assert.Equal(t, input.ID, result.ID)
			assert.Equal(t, input.DomainNames, result.DomainNames)
			assert.Equal(t, input.ProviderID, result.ProviderID)
			assert.Equal(t, input.IssuedAt, result.IssuedAt)
		})
	})

	t.Run("toIssueCertificateRequest", func(t *testing.T) {
		t.Run("converts issue request DTO to domain object", func(t *testing.T) {
			input := &issueCertificateRequest{
				ProviderID:  "digitalocean",
				DomainNames: []string{"example.com"},
				Parameters:  map[string]any{"token": "secret"},
			}

			result := toIssueCertificateRequest(input)

			assert.NotNil(t, result)
			assert.Equal(t, input.ProviderID, result.ProviderID)
			assert.Equal(t, input.DomainNames, result.DomainNames)
			assert.Equal(t, input.Parameters, result.Parameters)
		})
	})

	t.Run("toIssueCertificateResponse", func(t *testing.T) {
		t.Run("converts success case", func(t *testing.T) {
			id := uuid.New()
			cert := &certificate.Certificate{
				ID: id,
			}
			result := toIssueCertificateResponse(cert, nil)

			assert.True(t, result.Success)
			assert.Nil(t, result.ErrorReason)
			assert.Equal(t, &id, result.CertificateID)
		})

		t.Run("converts error case", func(t *testing.T) {
			err := errors.New("issue failed")
			result := toIssueCertificateResponse(nil, err)

			assert.False(t, result.Success)
			assert.Equal(t, "issue failed", *result.ErrorReason)
			assert.Nil(t, result.CertificateID)
		})
	})

	t.Run("toRenewCertificateResponse", func(t *testing.T) {
		t.Run("converts success case", func(t *testing.T) {
			result := toRenewCertificateResponse(nil)
			assert.True(t, result.Success)
			assert.Nil(t, result.ErrorReason)
		})

		t.Run("converts error case", func(t *testing.T) {
			err := errors.New("renew failed")
			result := toRenewCertificateResponse(err)
			assert.False(t, result.Success)
			assert.Equal(t, "renew failed", *result.ErrorReason)
		})
	})
}
