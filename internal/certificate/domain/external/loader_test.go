package external

import (
	"crypto/x509"
	"encoding/base64"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/certificate"
)

func Test_loader(t *testing.T) {
	ecKey := newECDSAKey(t)
	serverCertificate := newCertificate(t, ecKey, "example.com")

	caKey := newRSAKey(t)
	caCertificate := newCertificate(t, caKey, "ca.example.com")

	t.Run("loadCertificateFromParameters", func(t *testing.T) {
		t.Run("normalizes PEM files", func(t *testing.T) {
			directory := t.TempDir()

			privateKeyPath := filepath.Join(directory, "private-key.pem")
			require.NoError(t, os.WriteFile(privateKeyPath, encodeSEC1PEM(t, ecKey), 0o600))

			certificatePath := filepath.Join(directory, "certificate.pem")
			require.NoError(
				t,
				os.WriteFile(certificatePath, encodeCertificatePEM(serverCertificate), 0o600),
			)

			chainPath := filepath.Join(directory, "chain.pem")
			require.NoError(
				t,
				os.WriteFile(chainPath, encodeCertificateChainPEM(caCertificate), 0o600),
			)

			loaded, err := loadCertificateFromParameters(t.Context(), map[string]any{
				publicKeyPathFieldID:  certificatePath,
				privateKeyPathFieldID: privateKeyPath,
				chainPathFieldID:      chainPath,
			}, []string{"example.com"}, uuid.New())
			require.NoError(t, err)

			assertLoadedCertificate(t, loaded, serverCertificate, 1)
		})

		t.Run("normalizes DER files", func(t *testing.T) {
			rsaKey := newRSAKey(t)
			rsaCertificate := newCertificate(t, rsaKey, "example.com")

			directory := t.TempDir()

			privateKeyPath := filepath.Join(directory, "private-key.der")
			require.NoError(t, os.WriteFile(privateKeyPath, encodePKCS1DER(rsaKey), 0o600))

			certificatePath := filepath.Join(directory, "certificate.der")
			require.NoError(
				t,
				os.WriteFile(certificatePath, rsaCertificate.Raw, 0o600),
			)

			loaded, err := loadCertificateFromParameters(t.Context(), map[string]any{
				publicKeyPathFieldID:  certificatePath,
				privateKeyPathFieldID: privateKeyPath,
			}, []string{"example.com"}, uuid.New())
			require.NoError(t, err)

			assertLoadedCertificate(t, loaded, rsaCertificate, 0)
		})

		t.Run("returns error for invalid private key", func(t *testing.T) {
			directory := t.TempDir()

			privateKeyPath := filepath.Join(directory, "private-key.pem")
			require.NoError(t, os.WriteFile(privateKeyPath, []byte("invalid"), 0o600))

			certificatePath := filepath.Join(directory, "certificate.pem")
			require.NoError(
				t,
				os.WriteFile(certificatePath, encodeCertificatePEM(serverCertificate), 0o600),
			)

			_, err := loadCertificateFromParameters(t.Context(), map[string]any{
				publicKeyPathFieldID:  certificatePath,
				privateKeyPathFieldID: privateKeyPath,
			}, []string{"example.com"}, uuid.New())
			assert.Error(t, err)
		})

		t.Run("returns error for invalid certificate", func(t *testing.T) {
			directory := t.TempDir()

			privateKeyPath := filepath.Join(directory, "private-key.pem")
			require.NoError(t, os.WriteFile(privateKeyPath, encodeSEC1PEM(t, ecKey), 0o600))

			certificatePath := filepath.Join(directory, "certificate.pem")
			require.NoError(t, os.WriteFile(certificatePath, []byte("invalid"), 0o600))

			_, err := loadCertificateFromParameters(t.Context(), map[string]any{
				publicKeyPathFieldID:  certificatePath,
				privateKeyPathFieldID: privateKeyPath,
			}, []string{"example.com"}, uuid.New())
			assert.Error(t, err)
		})
	})
}

func assertLoadedCertificate(
	t *testing.T,
	loaded *certificate.Certificate,
	expectedCertificate *x509.Certificate,
	chainLength int,
) {
	t.Helper()

	privateKeyDER, err := base64.StdEncoding.DecodeString(loaded.PrivateKey)
	require.NoError(t, err)

	_, err = x509.ParsePKCS8PrivateKey(privateKeyDER)
	require.NoError(t, err)

	publicKeyDER, err := base64.StdEncoding.DecodeString(loaded.PublicKey)
	require.NoError(t, err)

	publicCertificate, err := x509.ParseCertificate(publicKeyDER)
	require.NoError(t, err)
	assert.Equal(t, expectedCertificate.Raw, publicCertificate.Raw)

	require.Len(t, loaded.CertificationChain, chainLength)

	for _, chainElement := range loaded.CertificationChain {
		chainDER, err := base64.StdEncoding.DecodeString(chainElement)
		require.NoError(t, err)

		_, err = x509.ParseCertificate(chainDER)
		require.NoError(t, err)
	}
}
