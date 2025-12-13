package custom

import (
	"context"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"strings"
	"time"

	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/certificate/commons"
	"dillmann.com.br/nginx-ignition/core/certificate/server"
	"dillmann.com.br/nginx-ignition/core/common/coreerror"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
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

func (p *Provider) DynamicFields() []*dynamicfields.DynamicField {
	return []*dynamicfields.DynamicField{
		&uploadModeField,
		&publicKeyTextField,
		&privateKeyTextField,
		&certificationChainTextField,
		&publicKeyFileField,
		&privateKeyFileField,
		&certificationChainFileField,
	}
}

func (p *Provider) Priority() int {
	return 2
}

func (p *Provider) Issue(_ context.Context, request *server.IssueRequest) (*server.Certificate, error) {
	if err := commons.Validate(request, validationRules{p.DynamicFields()}); err != nil {
		return nil, err
	}

	params := request.Parameters
	fileUploadMode := params[uploadModeField.ID] == fileUploadModeID

	var privateKeyStr, publicKeyStr, chainStr string
	var chainPresent bool

	if fileUploadMode {
		privateKeyStr, _ = params[privateKeyFileField.ID].(string)
		publicKeyStr, _ = params[publicKeyFileField.ID].(string)
		chainStr, chainPresent = params[certificationChainFileField.ID].(string)
	} else {
		privateKeyStr, _ = params[privateKeyTextField.ID].(string)
		publicKeyStr, _ = params[publicKeyTextField.ID].(string)
		chainStr, chainPresent = params[certificationChainTextField.ID].(string)
	}

	privateKey, err := parsePrivateKey(privateKeyStr, fileUploadMode)
	if err != nil {
		return nil, coreerror.New("Invalid private key", true)
	}

	publicKey, err := parseCertificate(publicKeyStr, fileUploadMode)
	if err != nil {
		return nil, coreerror.New("Invalid public key", true)
	}

	var chain []*x509.Certificate
	if chainPresent && chainStr != "" {
		chain, err = parseCertificateChain(chainStr, fileUploadMode)
		if err != nil {
			return nil, coreerror.New("Invalid certification chain", true)
		}
	}

	return &server.Certificate{
		ID:                 uuid.New(),
		DomainNames:        request.DomainNames,
		ProviderID:         p.ID(),
		IssuedAt:           time.Now(),
		ValidUntil:         publicKey.NotAfter,
		ValidFrom:          publicKey.NotBefore,
		RenewAfter:         nil,
		PrivateKey:         base64.StdEncoding.EncodeToString(privateKey),
		PublicKey:          base64.StdEncoding.EncodeToString(publicKey.Raw),
		CertificationChain: encodeChain(chain),
		Parameters:         request.Parameters,
		Metadata:           nil,
	}, nil
}

func (p *Provider) Renew(_ context.Context, cert *server.Certificate) (*server.Certificate, error) {
	return cert, nil
}

func parsePrivateKey(key string, base64Encoded bool) ([]byte, error) {
	decodedKey, err := stringToByteArray(key, base64Encoded)
	if err != nil {
		return nil, coreerror.New("Failed to decode key", true)
	}

	block, _ := pem.Decode(decodedKey)
	if block == nil {
		return nil, coreerror.New("Failed to parse PEM block containing the key", true)
	}

	return block.Bytes, nil
}

func parseCertificate(cert string, base64Encoded bool) (*x509.Certificate, error) {
	decodedCert, err := stringToByteArray(cert, base64Encoded)
	if err != nil {
		return nil, coreerror.New("Failed to decode certificate", true)
	}

	block, _ := pem.Decode(decodedCert)
	if block == nil {
		return nil, coreerror.New("Failed to parse PEM block containing the certificate", true)
	}

	return x509.ParseCertificate(block.Bytes)
}

func parseCertificateChain(chain string, base64Encoded bool) ([]*x509.Certificate, error) {
	decodedChain, err := stringToByteArray(chain, base64Encoded)
	if err != nil {
		return nil, coreerror.New("Failed to decode chain", true)
	}

	certs := make([]*x509.Certificate, 0)
	for _, cert := range strings.Split(string(decodedChain), "-----END CERTIFICATE-----") {
		if cert == "" {
			continue
		}

		cert += "-----END CERTIFICATE-----"
		parsedCert, err := parseCertificate(cert, false)
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

func stringToByteArray(value string, base64Encoded bool) ([]byte, error) {
	if base64Encoded {
		return base64.StdEncoding.DecodeString(value)
	}

	return []byte(value), nil
}
