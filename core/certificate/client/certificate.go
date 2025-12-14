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

	var certBytes []byte

	if ca {
		template.IsCA = true
		template.KeyUsage |= x509.KeyUsageCertSign
		certBytes, err = x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	} else {
		if caPrivateKey == nil {
			certBytes, err = x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
		} else {
			caKeyBlock, _ := pem.Decode(*caPrivateKey)
			if caKeyBlock == nil {
				return nil, err
			}

			var caPKey interface{}
			if caPKey, err = x509.ParsePKCS1PrivateKey(caKeyBlock.Bytes); err != nil {
				if caPKey, err = x509.ParsePKCS8PrivateKey(caKeyBlock.Bytes); err != nil {
					return nil, err
				}
			}

			certBytes, err = x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, caPKey)
		}
	}

	if err != nil {
		return nil, err
	}

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
	}, nil
}
