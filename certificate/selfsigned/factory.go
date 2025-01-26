package selfsigned

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"time"
)

const keySize = 2048

type certificateFactory struct {
	caPrivateKey  *rsa.PrivateKey
	caCertificate *x509.Certificate
}

func newFactory() (*certificateFactory, error) {
	caKeyPair, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return nil, err
	}

	caCert, err := buildCertificate(caKeyPair, caKeyPair, "nginx ignition", nil)
	if err != nil {
		return nil, err
	}

	return &certificateFactory{
		caPrivateKey:  caKeyPair,
		caCertificate: caCert,
	}, nil
}

func (f *certificateFactory) build(domainNames []string) (*x509.Certificate, *rsa.PrivateKey, error) {
	mainDomainName := domainNames[0]
	keyPair, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return nil, nil, err
	}

	altNames := domainNames[1:]
	cert, err := buildCertificate(f.caPrivateKey, keyPair, mainDomainName, altNames)
	if err != nil {
		return nil, nil, err
	}

	return cert, keyPair, nil
}

func buildCertificate(
	caPrivateKey *rsa.PrivateKey,
	keyPair *rsa.PrivateKey,
	commonName string,
	altNames []string,
) (*x509.Certificate, error) {
	notBefore := time.Now().Add(-1 * time.Hour)
	notAfter := time.Now().Add(365 * 24 * time.Hour)

	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return nil, err
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName: commonName,
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		BasicConstraintsValid: true,
	}

	if len(altNames) > 0 {
		template.DNSNames = altNames
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &keyPair.PublicKey, caPrivateKey)
	if err != nil {
		return nil, err
	}

	cert, err := x509.ParseCertificate(certBytes)
	if err != nil {
		return nil, err
	}

	return cert, nil
}
