package selfsigned

import (
	"crypto/x509"
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/certificate"
)

func Test_provider(t *testing.T) {
	t.Run("Issue", func(t *testing.T) {
		t.Run("stores PKCS8 private key", func(t *testing.T) {
			issued, err := New().Issue(t.Context(), &certificate.IssueRequest{
				DomainNames: []string{"example.com"},
				ProviderID:  "SELF_SIGNED",
				Parameters:  map[string]any{},
			})
			require.NoError(t, err)

			privateKeyDER, err := base64.StdEncoding.DecodeString(issued.PrivateKey)
			require.NoError(t, err)

			_, err = x509.ParsePKCS8PrivateKey(privateKeyDER)
			require.NoError(t, err)

			publicKeyDER, err := base64.StdEncoding.DecodeString(issued.PublicKey)
			require.NoError(t, err)

			_, err = x509.ParseCertificate(publicKeyDER)
			require.NoError(t, err)
		})
	})
}
