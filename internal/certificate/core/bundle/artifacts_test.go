package bundle

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func newRSAKey(t *testing.T) *rsa.PrivateKey {
	t.Helper()

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	return key
}

func newECDSAKey(t *testing.T) *ecdsa.PrivateKey {
	t.Helper()

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)

	return key
}

func newEd25519Key(t *testing.T) ed25519.PrivateKey {
	t.Helper()

	_, key, err := ed25519.GenerateKey(rand.Reader)
	require.NoError(t, err)

	return key
}

func newCertificate(t *testing.T, key crypto.Signer, domainNames ...string) *x509.Certificate {
	t.Helper()

	commonName := "example.com"
	if len(domainNames) > 0 {
		commonName = domainNames[0]
	}

	template := &x509.Certificate{
		SerialNumber:          big.NewInt(time.Now().UnixNano()),
		Subject:               pkix.Name{CommonName: commonName},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(24 * time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              domainNames,
	}

	der, err := x509.CreateCertificate(rand.Reader, template, template, key.Public(), key)
	require.NoError(t, err)

	certificate, err := x509.ParseCertificate(der)
	require.NoError(t, err)

	return certificate
}

func encodePKCS1PEM(key *rsa.PrivateKey) []byte {
	return encodePEM("RSA PRIVATE KEY", x509.MarshalPKCS1PrivateKey(key))
}

func encodePKCS1DER(key *rsa.PrivateKey) []byte {
	return x509.MarshalPKCS1PrivateKey(key)
}

func encodeSEC1PEM(t *testing.T, key *ecdsa.PrivateKey) []byte {
	t.Helper()

	return encodePEM("EC PRIVATE KEY", encodeSEC1DER(t, key))
}

func encodeSEC1DER(t *testing.T, key *ecdsa.PrivateKey) []byte {
	t.Helper()

	der, err := x509.MarshalECPrivateKey(key)
	require.NoError(t, err)

	return der
}

func encodePKCS8PEM(t *testing.T, key crypto.Signer) []byte {
	t.Helper()

	return encodePEM("PRIVATE KEY", encodePKCS8DER(t, key))
}

func encodePKCS8DER(t *testing.T, key crypto.Signer) []byte {
	t.Helper()

	der, err := x509.MarshalPKCS8PrivateKey(key)
	require.NoError(t, err)

	return der
}

func encodePEM(blockType string, der []byte) []byte {
	return pem.EncodeToMemory(&pem.Block{
		Type:  blockType,
		Bytes: der,
	})
}

func encodeCertificatePEM(certificate *x509.Certificate) []byte {
	return encodePEM("CERTIFICATE", certificate.Raw)
}

func encodeCertificateChainPEM(certificates ...*x509.Certificate) []byte {
	parts := make([][]byte, len(certificates))
	for index, certificate := range certificates {
		parts[index] = encodeCertificatePEM(certificate)
	}

	return bytes.Join(parts, nil)
}

func encodeCertificateChainDER(certificates ...*x509.Certificate) []byte {
	parts := make([][]byte, len(certificates))
	for index, certificate := range certificates {
		parts[index] = certificate.Raw
	}

	return bytes.Join(parts, nil)
}

func publicKeyDER(t *testing.T, key crypto.Signer) []byte {
	t.Helper()

	der, err := x509.MarshalPKIXPublicKey(key.Public())
	require.NoError(t, err)

	return der
}
