package bundle

import (
	"crypto"
	"crypto/ecdh"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/coreerror"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/i18n"
)

func Test_bundle(t *testing.T) {
	rsaKey := newRSAKey(t)
	ecKey := newECDSAKey(t)
	ed25519Key := newEd25519Key(t)
	certificate := newCertificate(t, ecKey, "example.com")

	t.Run("NormalizePrivateKey", func(t *testing.T) {
		t.Run("normalizes PKCS1 PEM", func(t *testing.T) {
			assertNormalizedPrivateKey(t, encodePKCS1PEM(rsaKey), rsaKey)
		})

		t.Run("normalizes PKCS1 DER", func(t *testing.T) {
			assertNormalizedPrivateKey(t, encodePKCS1DER(rsaKey), rsaKey)
		})

		t.Run("normalizes PKCS8 PEM", func(t *testing.T) {
			assertNormalizedPrivateKey(t, encodePKCS8PEM(t, rsaKey), rsaKey)
		})

		t.Run("normalizes PKCS8 DER", func(t *testing.T) {
			assertNormalizedPrivateKey(t, encodePKCS8DER(t, rsaKey), rsaKey)
		})

		t.Run("normalizes SEC1 PEM", func(t *testing.T) {
			assertNormalizedPrivateKey(t, encodeSEC1PEM(t, ecKey), ecKey)
		})

		t.Run("normalizes SEC1 DER", func(t *testing.T) {
			assertNormalizedPrivateKey(t, encodeSEC1DER(t, ecKey), ecKey)
		})

		t.Run("normalizes EC PKCS8 PEM", func(t *testing.T) {
			assertNormalizedPrivateKey(t, encodePKCS8PEM(t, ecKey), ecKey)
		})

		t.Run("normalizes Ed25519 PKCS8 PEM", func(t *testing.T) {
			assertNormalizedPrivateKey(t, encodePKCS8PEM(t, ed25519Key), ed25519Key)
		})

		t.Run("skips EC PARAMETERS before key", func(t *testing.T) {
			raw := append(
				encodePEM("EC PARAMETERS", []byte("parameters")),
				encodeSEC1PEM(t, ecKey)...,
			)

			assertNormalizedPrivateKey(t, raw, ecKey)
		})

		t.Run("skips certificate before key", func(t *testing.T) {
			raw := append(encodeCertificatePEM(certificate), encodeSEC1PEM(t, ecKey)...)

			assertNormalizedPrivateKey(t, raw, ecKey)
		})

		t.Run("normalizes mislabeled PKCS1 PEM", func(t *testing.T) {
			raw := encodePEM("PRIVATE KEY", encodePKCS1DER(rsaKey))

			assertNormalizedPrivateKey(t, raw, rsaKey)
		})

		t.Run("normalizes mislabeled SEC1 PEM", func(t *testing.T) {
			raw := encodePEM("PRIVATE KEY", encodeSEC1DER(t, ecKey))

			assertNormalizedPrivateKey(t, raw, ecKey)
		})

		t.Run("skips leading text", func(t *testing.T) {
			raw := append([]byte("leading text\n"), encodeSEC1PEM(t, ecKey)...)

			assertNormalizedPrivateKey(t, raw, ecKey)
		})
	})

	t.Run("ParsePrivateKey", func(t *testing.T) {
		t.Run("rejects empty input", func(t *testing.T) {
			_, err := ParsePrivateKey(t.Context(), nil)
			assertBundleError(t, err, i18n.K.CertificateBundleInvalidPrivateKey)
		})

		t.Run("rejects invalid input", func(t *testing.T) {
			_, err := ParsePrivateKey(t.Context(), []byte("not a key"))
			assertBundleError(t, err, i18n.K.CertificateBundleInvalidPrivateKey)
		})

		t.Run("rejects encrypted PKCS8", func(t *testing.T) {
			raw := encodePEM("ENCRYPTED PRIVATE KEY", []byte("encrypted"))

			_, err := ParsePrivateKey(t.Context(), raw)
			assertBundleError(t, err, i18n.K.CertificateBundleEncryptedPrivateKey)
		})

		t.Run("rejects legacy encrypted PEM", func(t *testing.T) {
			raw := pem.EncodeToMemory(&pem.Block{
				Type:    "RSA PRIVATE KEY",
				Headers: map[string]string{"Proc-Type": "4,ENCRYPTED"},
				Bytes:   encodePKCS1DER(rsaKey),
			})

			_, err := ParsePrivateKey(t.Context(), raw)
			assertBundleError(t, err, i18n.K.CertificateBundleEncryptedPrivateKey)
		})

		t.Run("rejects OpenSSH key", func(t *testing.T) {
			raw := encodePEM("OPENSSH PRIVATE KEY", []byte("key"))

			_, err := ParsePrivateKey(t.Context(), raw)
			assertBundleError(t, err, i18n.K.CertificateBundleUnsupportedPrivateKey)
		})

		t.Run("rejects DSA key", func(t *testing.T) {
			raw := encodePEM("DSA PRIVATE KEY", []byte("key"))

			_, err := ParsePrivateKey(t.Context(), raw)
			assertBundleError(t, err, i18n.K.CertificateBundleUnsupportedPrivateKey)
		})

		t.Run("rejects X25519 key", func(t *testing.T) {
			key, err := ecdh.X25519().GenerateKey(rand.Reader)
			require.NoError(t, err)

			der, err := x509.MarshalPKCS8PrivateKey(key)
			require.NoError(t, err)

			_, err = ParsePrivateKey(t.Context(), encodePEM("PRIVATE KEY", der))
			assertBundleError(t, err, i18n.K.CertificateBundleUnsupportedPrivateKey)
		})
	})

	t.Run("ParseCertificate", func(t *testing.T) {
		trustedDER := append(append([]byte{}, certificate.Raw...), 0x30, 0x00)

		t.Run("parses PEM certificate", func(t *testing.T) {
			parsed, err := ParseCertificate(t.Context(), encodeCertificatePEM(certificate))
			require.NoError(t, err)
			assert.Equal(t, certificate.Raw, parsed.Raw)
		})

		t.Run("parses DER certificate", func(t *testing.T) {
			parsed, err := ParseCertificate(t.Context(), certificate.Raw)
			require.NoError(t, err)
			assert.Equal(t, certificate.Raw, parsed.Raw)
		})

		t.Run("parses X509 CERTIFICATE block", func(t *testing.T) {
			raw := encodePEM("X509 CERTIFICATE", certificate.Raw)

			parsed, err := ParseCertificate(t.Context(), raw)
			require.NoError(t, err)
			assert.Equal(t, certificate.Raw, parsed.Raw)
		})

		t.Run("parses TRUSTED CERTIFICATE block", func(t *testing.T) {
			raw := encodePEM("TRUSTED CERTIFICATE", trustedDER)

			parsed, err := ParseCertificate(t.Context(), raw)
			require.NoError(t, err)
			assert.Equal(t, certificate.Raw, parsed.Raw)
		})

		t.Run("skips OpenSSL bag attributes", func(t *testing.T) {
			raw := append(
				[]byte("Bag Attributes\n    friendlyName: test\n"),
				encodeCertificatePEM(certificate)...,
			)

			parsed, err := ParseCertificate(t.Context(), raw)
			require.NoError(t, err)
			assert.Equal(t, certificate.Raw, parsed.Raw)
		})

		t.Run("skips key before certificate", func(t *testing.T) {
			raw := append(encodeSEC1PEM(t, ecKey), encodeCertificatePEM(certificate)...)

			parsed, err := ParseCertificate(t.Context(), raw)
			require.NoError(t, err)
			assert.Equal(t, certificate.Raw, parsed.Raw)
		})

		t.Run("returns the first certificate from a chain", func(t *testing.T) {
			second := newCertificate(t, rsaKey, "ca.example.com")
			raw := encodeCertificateChainPEM(certificate, second)

			parsed, err := ParseCertificate(t.Context(), raw)
			require.NoError(t, err)
			assert.Equal(t, certificate.Raw, parsed.Raw)
		})
	})

	t.Run("ParseCertificateChain", func(t *testing.T) {
		second := newCertificate(t, rsaKey, "ca.example.com")

		t.Run("parses PEM bundle", func(t *testing.T) {
			parsed, err := ParseCertificateChain(
				t.Context(),
				encodeCertificateChainPEM(certificate, second),
			)
			require.NoError(t, err)
			require.Len(t, parsed, 2)
			assert.Equal(t, certificate.Raw, parsed[0].Raw)
			assert.Equal(t, second.Raw, parsed[1].Raw)
		})

		t.Run("parses concatenated DER", func(t *testing.T) {
			parsed, err := ParseCertificateChain(
				t.Context(),
				encodeCertificateChainDER(certificate, second),
			)
			require.NoError(t, err)
			require.Len(t, parsed, 2)
			assert.Equal(t, certificate.Raw, parsed[0].Raw)
			assert.Equal(t, second.Raw, parsed[1].Raw)
		})

		t.Run("rejects empty input", func(t *testing.T) {
			_, err := ParseCertificateChain(t.Context(), nil)
			assertBundleError(t, err, i18n.K.CertificateBundleInvalidCertificate)
		})

		t.Run("rejects invalid input", func(t *testing.T) {
			_, err := ParseCertificateChain(t.Context(), []byte("not a certificate"))
			assertBundleError(t, err, i18n.K.CertificateBundleInvalidCertificate)
		})

		t.Run("rejects PEM without certificates", func(t *testing.T) {
			_, err := ParseCertificateChain(t.Context(), encodeSEC1PEM(t, ecKey))
			assertBundleError(t, err, i18n.K.CertificateBundleInvalidCertificate)
		})

		t.Run("rejects malformed certificate", func(t *testing.T) {
			_, err := ParseCertificateChain(
				t.Context(),
				encodePEM("CERTIFICATE", []byte("invalid")),
			)
			assertBundleError(t, err, i18n.K.CertificateBundleInvalidCertificate)
		})
	})
}

func assertNormalizedPrivateKey(t *testing.T, raw []byte, expectedKey crypto.Signer) {
	t.Helper()

	normalized, err := NormalizePrivateKey(t.Context(), raw)
	require.NoError(t, err)

	block, _ := pem.Decode(normalized)
	assert.Nil(t, block)

	parsed, err := x509.ParsePKCS8PrivateKey(normalized)
	require.NoError(t, err)

	signer, ok := parsed.(crypto.Signer)
	require.True(t, ok)
	assert.Equal(t, publicKeyDER(t, expectedKey), publicKeyDER(t, signer))
}

func assertBundleError(t *testing.T, err error, expectedKey string) {
	t.Helper()

	coreErr := &coreerror.CoreError{}
	require.ErrorAs(t, err, &coreErr)
	assert.Equal(t, expectedKey, coreErr.Message.Key)
}
