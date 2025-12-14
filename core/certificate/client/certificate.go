package client

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"time"
)

type CertificateRequest struct {
	Name             string
	CommonName       string
	Email            string
	Organization     string
	OrganizationUnit string
	Validity         time.Duration
	KeySize          int
}

type CertificateResponse struct {
	PrivateKey *[]byte
	PublicKey  *[]byte
	IssuedAt   time.Time
	ExpiresAt  time.Time
}

func buildCertificate(request *CertificateRequest, ca bool, caPrivateKey *[]byte) (*CertificateResponse, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, request.KeySize)
	if err != nil {
		return nil, err
	}

	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return nil, err
	}

	template := buildCertificateTemplate(request, serialNumber)

	certBytes, err := createCertificateBytes(&template, privateKey, ca, caPrivateKey)
	if err != nil {
		return nil, err
	}

	return buildCertificateResponse(privateKey, certBytes, &template), nil
}

func buildCertificateTemplate(request *CertificateRequest, serialNumber *big.Int) x509.Certificate {
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName:         request.CommonName,
			Organization:       []string{request.Organization},
			OrganizationalUnit: []string{request.OrganizationUnit},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(request.Validity),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
	}

	if request.Email != "" {
		template.EmailAddresses = []string{request.Email}
	}

	return template
}

func createCertificateBytes(
	template *x509.Certificate,
	privateKey *rsa.PrivateKey,
	ca bool,
	caPrivateKey *[]byte,
) ([]byte, error) {
	if ca {
		template.IsCA = true
		template.KeyUsage |= x509.KeyUsageCertSign
		return x509.CreateCertificate(rand.Reader, template, template, &privateKey.PublicKey, privateKey)
	}

	if caPrivateKey == nil {
		return x509.CreateCertificate(rand.Reader, template, template, &privateKey.PublicKey, privateKey)
	}

	caPKey, err := parseCAPrivateKey(caPrivateKey)
	if err != nil {
		return nil, err
	}

	return x509.CreateCertificate(rand.Reader, template, template, &privateKey.PublicKey, caPKey)
}

func parseCAPrivateKey(caPrivateKey *[]byte) (interface{}, error) {
	caKeyBlock, _ := pem.Decode(*caPrivateKey)
	if caKeyBlock == nil {
		return nil, x509.ErrUnsupportedAlgorithm
	}

	caPKey, err := x509.ParsePKCS1PrivateKey(caKeyBlock.Bytes)
	if err == nil {
		return caPKey, nil
	}

	return x509.ParsePKCS8PrivateKey(caKeyBlock.Bytes)
}

func buildCertificateResponse(
	privateKey *rsa.PrivateKey,
	certBytes []byte,
	template *x509.Certificate,
) *CertificateResponse {
	privateKeyBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	publicKeyBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	return &CertificateResponse{
		PrivateKey: &privateKeyBytes,
		PublicKey:  &publicKeyBytes,
		IssuedAt:   template.NotBefore,
		ExpiresAt:  template.NotAfter,
	}
}
