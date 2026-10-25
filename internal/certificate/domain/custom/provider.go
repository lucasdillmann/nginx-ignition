package custom

import (
	"context"
	"crypto/x509"
	"encoding/base64"
	"time"

	"github.com/google/uuid"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/coreerror"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/dynamicfields"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/i18n"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/certificate"
	"github.com/lucasdillmann/nginx-ignition/internal/certificate/core/bundle"
	"github.com/lucasdillmann/nginx-ignition/internal/certificate/core/validator"
)

type Provider struct{}

func New() *Provider {
	return &Provider{}
}

func (p *Provider) ID() string {
	return "CUSTOM"
}

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateCustomName)
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
	if err := validator.Validate(ctx, request, validationRules{p.DynamicFields(ctx)}); err != nil {
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
		return nil, err
	}

	publicKey, err := parseCertificate(ctx, publicKeyStr, fileUploadMode)
	if err != nil {
		return nil, err
	}

	chain := make([]*x509.Certificate, 0)
	if chainPresent && chainStr != "" {
		chain, err = parseCertificateChain(ctx, chainStr, fileUploadMode)
		if err != nil {
			return nil, err
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

func (p *Provider) IsDueToRenew(_ context.Context, _ *certificate.Certificate) (bool, error) {
	return false, nil
}

func parsePrivateKey(ctx context.Context, key string, base64Encoded bool) ([]byte, error) {
	decodedKey, err := stringToByteArray(key, base64Encoded)
	if err != nil {
		return nil, coreerror.New(
			i18n.M(ctx, i18n.K.CertificateCustomUnableToDecode).V("type", "key"),
			true,
		)
	}

	return bundle.NormalizePrivateKey(ctx, decodedKey)
}

func parseCertificate(
	ctx context.Context,
	cert string,
	base64Encoded bool,
) (*x509.Certificate, error) {
	decodedCert, err := stringToByteArray(cert, base64Encoded)
	if err != nil {
		return nil, coreerror.New(
			i18n.M(ctx, i18n.K.CertificateCustomUnableToDecode).V("type", "certificate"),
			true,
		)
	}

	return bundle.ParseCertificate(ctx, decodedCert)
}

func parseCertificateChain(
	ctx context.Context,
	chain string,
	base64Encoded bool,
) ([]*x509.Certificate, error) {
	decodedChain, err := stringToByteArray(chain, base64Encoded)
	if err != nil {
		return nil, coreerror.New(
			i18n.M(ctx, i18n.K.CertificateCustomUnableToDecode).V("type", "chain"),
			true,
		)
	}

	return bundle.ParseCertificateChain(ctx, decodedChain)
}

func encodeChain(chain []*x509.Certificate) []string {
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
