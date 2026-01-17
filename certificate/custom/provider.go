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
	"dillmann.com.br/nginx-ignition/core/certificate"
	"dillmann.com.br/nginx-ignition/core/common/coreerror"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

type Provider struct{}

func New() *Provider {
	return &Provider{}
}

func (p *Provider) ID() string {
	return "CUSTOM"
}

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateCommonCustomName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dynamicFields(ctx)
}

func (p *Provider) Priority() int {
	return 2
}

func (p *Provider) Issue(
	ctx context.Context,
	request *certificate.IssueRequest,
) (*certificate.Certificate, error) {
	if err := commons.Validate(ctx, request, validationRules{p.DynamicFields(ctx)}); err != nil {
		return nil, err
	}

	params := request.Parameters
	fileUploadMode := params[uploadModeFieldID] == fileUploadModeID

	var privateKeyStr, publicKeyStr, chainStr string
	var chainPresent bool

	if fileUploadMode {
		privateKeyStr, _ = params[privateKeyFileFieldID].(string)
		publicKeyStr, _ = params[publicKeyFileFieldID].(string)
		chainStr, chainPresent = params[chainFileFieldID].(string)
	} else {
		privateKeyStr, _ = params[privateKeyTextFieldID].(string)
		publicKeyStr, _ = params[publicKeyTextFieldID].(string)
		chainStr, chainPresent = params[chainTextFieldID].(string)
	}

	privateKey, err := parsePrivateKey(ctx, privateKeyStr, fileUploadMode)
	if err != nil {
		return nil, coreerror.New(i18n.M(ctx, i18n.K.CertificateErrorInvalidPrivateKey), true)
	}

	publicKey, err := parseCertificate(ctx, publicKeyStr, fileUploadMode)
	if err != nil {
		return nil, coreerror.New(i18n.M(ctx, i18n.K.CertificateErrorInvalidPublicKey), true)
	}

	chain := make([]x509.Certificate, 0)
	if chainPresent && chainStr != "" {
		chain, err = parseCertificateChain(ctx, chainStr, fileUploadMode)
		if err != nil {
			return nil, coreerror.New(
				i18n.M(ctx, i18n.K.CertificateErrorInvalidCertificationChain),
				true,
			)
		}
	}

	return &certificate.Certificate{
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

func (p *Provider) Renew(
	_ context.Context,
	cert *certificate.Certificate,
) (*certificate.Certificate, error) {
	return cert, nil
}

func parsePrivateKey(ctx context.Context, key string, base64Encoded bool) ([]byte, error) {
	decodedKey, err := stringToByteArray(key, base64Encoded)
	if err != nil {
		return nil, coreerror.New(
			i18n.M(ctx, i18n.K.CertificateErrorUnableToDecode).V("type", "key"),
			true,
		)
	}

	block, _ := pem.Decode(decodedKey)
	if block == nil {
		return nil, coreerror.New(
			i18n.M(ctx, i18n.K.CertificateErrorUnableToParsePem).V("type", "key"),
			true,
		)
	}

	return block.Bytes, nil
}

func parseCertificate(
	ctx context.Context,
	cert string,
	base64Encoded bool,
) (*x509.Certificate, error) {
	decodedCert, err := stringToByteArray(cert, base64Encoded)
	if err != nil {
		return nil, coreerror.New(
			i18n.M(ctx, i18n.K.CertificateErrorUnableToDecode).V("type", "certificate"),
			true,
		)
	}

	block, _ := pem.Decode(decodedCert)
	if block == nil {
		return nil, coreerror.New(
			i18n.M(ctx, i18n.K.CertificateErrorUnableToParsePem).V("type", "certificate"),
			true,
		)
	}

	return x509.ParseCertificate(block.Bytes)
}

func parseCertificateChain(
	ctx context.Context,
	chain string,
	base64Encoded bool,
) ([]x509.Certificate, error) {
	decodedChain, err := stringToByteArray(chain, base64Encoded)
	if err != nil {
		return nil, coreerror.New(
			i18n.M(ctx, i18n.K.CertificateErrorUnableToDecode).V("type", "chain"),
			true,
		)
	}

	certs := make([]x509.Certificate, 0)
	for _, cert := range strings.Split(string(decodedChain), "-----END CERTIFICATE-----") {
		if cert == "" {
			continue
		}

		cert += "-----END CERTIFICATE-----"
		parsedCert, err := parseCertificate(ctx, cert, false)
		if err != nil {
			return nil, err
		}

		certs = append(certs, *parsedCert)
	}

	return certs, nil
}

func encodeChain(chain []x509.Certificate) []string {
	encodedChain := make([]string, len(chain))
	for index, cert := range chain {
		encodedChain[index] = base64.StdEncoding.EncodeToString(cert.Raw)
	}

	return encodedChain
}

func stringToByteArray(value string, base64Encoded bool) ([]byte, error) {
	if base64Encoded {
		return base64.StdEncoding.DecodeString(value)
	}

	return []byte(value), nil
}
