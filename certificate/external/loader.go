package external

import (
	"context"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"os"
	"time"

	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/certificate"
	"dillmann.com.br/nginx-ignition/core/common/coreerror"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

func loadCertificateFromParameters(
	ctx context.Context,
	params map[string]any,
	domainNames []string,
	certificateID uuid.UUID,
) (*certificate.Certificate, error) {
	publicPath, _ := params[publicKeyPathFieldID].(string)
	privatePath, _ := params[privateKeyPathFieldID].(string)
	chainPath, _ := params[chainPathFieldID].(string)

	privateKeyBytes, err := readPrivateKeyPEMBytes(ctx, privatePath)
	if err != nil {
		return nil, coreerror.New(i18n.M(ctx, i18n.K.CertificateExternalInvalidPrivateKey), true)
	}

	publicCert, err := readCertificate(ctx, publicPath)
	if err != nil {
		return nil, coreerror.New(i18n.M(ctx, i18n.K.CertificateExternalInvalidPublicKey), true)
	}

	chain := make([]x509.Certificate, 0)
	if chainPath != "" {
		chain, err = readCertificateChain(ctx, chainPath)
		if err != nil {
			return nil, coreerror.New(
				i18n.M(ctx, i18n.K.CertificateExternalInvalidCertificationChain),
				true,
			)
		}
	}

	return &certificate.Certificate{
		ID:                 certificateID,
		DomainNames:        domainNames,
		ProviderID:         providerID,
		IssuedAt:           time.Now(),
		ValidUntil:         publicCert.NotAfter,
		ValidFrom:          publicCert.NotBefore,
		RenewAfter:         nil,
		PrivateKey:         base64.StdEncoding.EncodeToString(privateKeyBytes),
		PublicKey:          base64.StdEncoding.EncodeToString(publicCert.Raw),
		CertificationChain: encodeChain(chain),
		Parameters:         params,
		Metadata:           nil,
	}, nil
}

func readPrivateKeyPEMBytes(ctx context.Context, path string) ([]byte, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, coreerror.New(
			i18n.M(ctx, i18n.K.CertificateExternalPathNotReadable).V("path", path),
			true,
		)
	}

	block, _ := pem.Decode(raw)
	if block == nil {
		return nil, coreerror.New(
			i18n.M(ctx, i18n.K.CommonUnableToParsePem).V("type", "key"),
			true,
		)
	}

	return block.Bytes, nil
}

func readCertificate(ctx context.Context, path string) (*x509.Certificate, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, coreerror.New(
			i18n.M(ctx, i18n.K.CertificateExternalPathNotReadable).V("path", path),
			true,
		)
	}

	block, _ := pem.Decode(raw)
	if block == nil {
		return nil, coreerror.New(
			i18n.M(ctx, i18n.K.CommonUnableToParsePem).V("type", "certificate"),
			true,
		)
	}

	return x509.ParseCertificate(block.Bytes)
}

func readCertificateChain(ctx context.Context, path string) ([]x509.Certificate, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, coreerror.New(
			i18n.M(ctx, i18n.K.CertificateExternalPathNotReadable).V("path", path),
			true,
		)
	}

	certificates := make([]x509.Certificate, 0)
	remaining := raw

	for {
		block, rest := pem.Decode(remaining)
		if block == nil {
			break
		}

		remaining = rest
		if block.Type != "CERTIFICATE" {
			continue
		}

		parsed, parseErr := x509.ParseCertificate(block.Bytes)
		if parseErr != nil {
			return nil, coreerror.New(
				i18n.M(ctx, i18n.K.CertificateExternalInvalidCertificationChain),
				true,
			)
		}

		certificates = append(certificates, *parsed)
	}

	if len(certificates) == 0 {
		return nil, coreerror.New(
			i18n.M(ctx, i18n.K.CertificateExternalInvalidCertificationChain),
			true,
		)
	}

	return certificates, nil
}

func encodeChain(chain []x509.Certificate) []string {
	result := make([]string, len(chain))
	for index, item := range chain {
		result[index] = base64.StdEncoding.EncodeToString(item.Raw)
	}

	return result
}
