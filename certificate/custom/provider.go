package custom

import (
	"crypto/x509"
	"dillmann.com.br/nginx-ignition/core/certificate"
	"dillmann.com.br/nginx-ignition/core/common/core_error"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
	"encoding/base64"
	"encoding/pem"
	"github.com/google/uuid"
	"strings"
	"time"
)

type Provider struct{}

func New() *Provider {
	return &Provider{}
}

func (p *Provider) ID() string {
	return "CUSTOM"
}

func (p *Provider) Name() string {
	return "Custom certificate"
}

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return []*dynamic_fields.DynamicField{
		&publicKeyField,
		&privateKeyField,
		&certificationChainField,
	}
}

func (p *Provider) Priority() int {
	return 2
}

func (p *Provider) Issue(request *certificate.IssueRequest) (*certificate.Certificate, error) {
	privateKeyStr := *request.Parameters[privateKeyField.ID]
	publicKeyStr := *request.Parameters[publicKeyField.ID]
	chainStr, chainPresent := request.Parameters[certificationChainField.ID]

	privateKey, err := parsePrivateKey(privateKeyStr.(string))
	if err != nil {
		return nil, core_error.New("Invalid private key", true)
	}

	publicKey, err := parseCertificate(publicKeyStr.(string))
	if err != nil {
		return nil, core_error.New("Invalid public key", true)
	}

	var chain []*x509.Certificate
	if chainPresent && chainStr != nil {
		chain, err = parseCertificateChain((*chainStr).(string))
		if err != nil {
			return nil, core_error.New("Invalid certification chain", true)
		}
	}

	return &certificate.Certificate{
		ID:                 uuid.New(),
		DomainNames:        dereferenceSlice(request.DomainNames),
		ProviderID:         p.ID(),
		IssuedAt:           time.Now(),
		ValidUntil:         publicKey.NotAfter,
		ValidFrom:          publicKey.NotBefore,
		RenewAfter:         nil,
		PrivateKey:         base64.StdEncoding.EncodeToString(privateKey),
		PublicKey:          base64.StdEncoding.EncodeToString(publicKey.Raw),
		CertificationChain: encodeChain(chain),
		Parameters:         dereferenceMap(request.Parameters),
		Metadata:           nil,
	}, nil
}

func (p *Provider) Renew(cert *certificate.Certificate) (*certificate.Certificate, error) {
	return cert, nil
}

func parsePrivateKey(key string) ([]byte, error) {
	decodedKey, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, core_error.New("Failed to decode key", true)
	}

	block, _ := pem.Decode(decodedKey)
	if block == nil {
		return nil, core_error.New("Failed to parse PEM block containing the key", true)
	}

	return block.Bytes, nil
}

func parseCertificate(cert string) (*x509.Certificate, error) {
	decodedCert, err := base64.StdEncoding.DecodeString(cert)
	if err != nil {
		return nil, core_error.New("Failed to decode certificate", true)
	}

	block, _ := pem.Decode(decodedCert)
	if block == nil {
		return nil, core_error.New("Failed to parse PEM block containing the certificate", true)
	}

	return x509.ParseCertificate(block.Bytes)
}

func parseCertificateChain(chain string) ([]*x509.Certificate, error) {
	decodedChain, err := base64.StdEncoding.DecodeString(chain)
	if err != nil {
		return nil, core_error.New("Failed to decode chain", true)
	}

	var certs []*x509.Certificate
	for _, cert := range strings.Split(string(decodedChain), "-----END CERTIFICATE-----") {
		if cert == "" {
			continue
		}

		cert += "-----END CERTIFICATE-----"
		parsedCert, err := parseCertificate(cert)
		if err != nil {
			return nil, err
		}

		certs = append(certs, parsedCert)
	}

	return certs, nil
}

func encodeChain(chain []*x509.Certificate) []string {
	encodedChain := make([]string, len(chain))
	for _, cert := range chain {
		encodedChain = append(encodedChain, base64.StdEncoding.EncodeToString(cert.Raw))
	}

	return encodedChain
}

func dereferenceSlice(input []*string) []string {
	output := make([]string, len(input))
	for index, value := range input {
		output[index] = *value
	}

	return output
}

func dereferenceMap(input map[string]*interface{}) map[string]interface{} {
	output := make(map[string]interface{})
	for key, value := range input {
		if value != nil {
			output[key] = *value
		}
	}

	return output
}
