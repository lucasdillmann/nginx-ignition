package custom

import (
	"crypto/x509"
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/certificate"
)

func Test_provider(t *testing.T) {
	ecKey := newECDSAKey(t)
	serverCertificate := newCertificate(t, ecKey, "example.com")

	caKey := newRSAKey(t)
	caCertificate := newCertificate(t, caKey, "ca.example.com")

	privateKeyPEM := encodeSEC1PEM(t, ecKey)
	certificatePEM := encodeCertificatePEM(serverCertificate)
	chainPEM := encodeCertificateChainPEM(caCertificate)

	t.Run("Issue", func(t *testing.T) {
		t.Run("normalizes text mode input", func(t *testing.T) {
			issued, err := New().Issue(t.Context(), &certificate.IssueRequest{
				DomainNames: []string{"example.com"},
				ProviderID:  "CUSTOM",
				Parameters: map[string]any{
					uploadModeFieldID:     textFieldUploadModeID,
					privateKeyTextFieldID: string(privateKeyPEM),
					publicKeyTextFieldID:  string(certificatePEM),
					chainTextFieldID:      string(chainPEM),
				},
			})
			require.NoError(t, err)

			assertNormalizedCertificate(t, issued, serverCertificate, 1)
		})

		t.Run("normalizes file mode input", func(t *testing.T) {
			issued, err := New().Issue(t.Context(), &certificate.IssueRequest{
				DomainNames: []string{"example.com"},
				ProviderID:  "CUSTOM",
				Parameters: map[string]any{
					uploadModeFieldID:     fileUploadModeID,
					privateKeyFileFieldID: base64.StdEncoding.EncodeToString(privateKeyPEM),
					publicKeyFileFieldID:  base64.StdEncoding.EncodeToString(certificatePEM),
					chainFileFieldID:      base64.StdEncoding.EncodeToString(chainPEM),
				},
			})
			require.NoError(t, err)

			assertNormalizedCertificate(t, issued, serverCertificate, 1)
		})

		t.Run("normalizes EC parameters before key", func(t *testing.T) {
			keyWithParameters := append(
				encodePEM("EC PARAMETERS", []byte("parameters")),
				privateKeyPEM...,
			)

			issued, err := New().Issue(t.Context(), &certificate.IssueRequest{
				DomainNames: []string{"example.com"},
				ProviderID:  "CUSTOM",
				Parameters: map[string]any{
					uploadModeFieldID:     textFieldUploadModeID,
					privateKeyTextFieldID: string(keyWithParameters),
					publicKeyTextFieldID:  string(certificatePEM),
				},
			})
			require.NoError(t, err)

			assertNormalizedCertificate(t, issued, serverCertificate, 0)
		})

		t.Run("normalizes RSA PKCS1 key", func(t *testing.T) {
			rsaKey := newRSAKey(t)
			rsaCertificate := newCertificate(t, rsaKey, "example.com")

			issued, err := New().Issue(t.Context(), &certificate.IssueRequest{
				DomainNames: []string{"example.com"},
				ProviderID:  "CUSTOM",
				Parameters: map[string]any{
					uploadModeFieldID:     textFieldUploadModeID,
					privateKeyTextFieldID: string(encodePKCS1PEM(rsaKey)),
					publicKeyTextFieldID:  string(encodeCertificatePEM(rsaCertificate)),
				},
			})
			require.NoError(t, err)

			assertNormalizedCertificate(t, issued, rsaCertificate, 0)
		})

		t.Run("returns error for invalid private key", func(t *testing.T) {
			_, err := New().Issue(t.Context(), &certificate.IssueRequest{
				DomainNames: []string{"example.com"},
				ProviderID:  "CUSTOM",
				Parameters: map[string]any{
					uploadModeFieldID:     textFieldUploadModeID,
					privateKeyTextFieldID: "not a key",
					publicKeyTextFieldID:  string(certificatePEM),
				},
			})
			assert.Error(t, err)
		})
	})
}

func assertNormalizedCertificate(
	t *testing.T,
	issued *certificate.Certificate,
	expectedCertificate *x509.Certificate,
	chainLength int,
) {
	t.Helper()

	privateKeyDER, err := base64.StdEncoding.DecodeString(issued.PrivateKey)
	require.NoError(t, err)

	_, err = x509.ParsePKCS8PrivateKey(privateKeyDER)
	require.NoError(t, err)

	publicKeyDER, err := base64.StdEncoding.DecodeString(issued.PublicKey)
	require.NoError(t, err)

	publicCertificate, err := x509.ParseCertificate(publicKeyDER)
	require.NoError(t, err)
	assert.Equal(t, expectedCertificate.Raw, publicCertificate.Raw)

	require.Len(t, issued.CertificationChain, chainLength)

	for _, chainElement := range issued.CertificationChain {
		chainDER, err := base64.StdEncoding.DecodeString(chainElement)
		require.NoError(t, err)

		_, err = x509.ParseCertificate(chainDER)
		require.NoError(t, err)
	}
}
