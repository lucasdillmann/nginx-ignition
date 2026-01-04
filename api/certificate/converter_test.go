package certificate

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_toCertificateResponse(t *testing.T) {
	t.Run("converts certificate to response DTO", func(t *testing.T) {
		input := newCertificate()
		result := toCertificateResponse(input)

		assert.NotNil(t, result)
		assert.Equal(t, input.ID, result.ID)
		assert.Equal(t, input.DomainNames, result.DomainNames)
		assert.Equal(t, input.ProviderID, result.ProviderID)
		assert.Equal(t, input.IssuedAt, result.IssuedAt)
	})
}

func Test_toIssueCertificateRequest(t *testing.T) {
	t.Run("converts issue request DTO to domain object", func(t *testing.T) {
		input := newIssueCertificateRequest()
		result := toIssueCertificateRequest(input)

		assert.NotNil(t, result)
		assert.Equal(t, input.ProviderID, result.ProviderID)
		assert.Equal(t, input.DomainNames, result.DomainNames)
		assert.Equal(t, input.Parameters, result.Parameters)
	})
}

func Test_toIssueCertificateResponse(t *testing.T) {
	t.Run("converts success case", func(t *testing.T) {
		certificateData := newCertificate()
		result := toIssueCertificateResponse(certificateData, nil)

		assert.True(t, result.Success)
		assert.Nil(t, result.ErrorReason)
		assert.Equal(t, &certificateData.ID, result.CertificateID)
	})

	t.Run("converts error case", func(t *testing.T) {
		expectedErr := errors.New("issue failed")
		result := toIssueCertificateResponse(nil, expectedErr)

		assert.False(t, result.Success)
		assert.Equal(t, "issue failed", *result.ErrorReason)
		assert.Nil(t, result.CertificateID)
	})
}

func Test_toRenewCertificateResponse(t *testing.T) {
	t.Run("converts success case", func(t *testing.T) {
		result := toRenewCertificateResponse(nil)
		assert.True(t, result.Success)
		assert.Nil(t, result.ErrorReason)
	})

	t.Run("converts error case", func(t *testing.T) {
		expectedErr := errors.New("renew failed")
		result := toRenewCertificateResponse(expectedErr)
		assert.False(t, result.Success)
		assert.Equal(t, "renew failed", *result.ErrorReason)
	})
}
