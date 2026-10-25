package external

import (
	"context"
	"crypto/x509"
	"encoding/base64"
	"os"
	"time"

	"github.com/google/uuid"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/coreerror"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/i18n"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/certificate"
	"github.com/lucasdillmann/nginx-ignition/internal/certificate/core/bundle"
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
		return nil, err
	}

	publicCert, err := readCertificate(ctx, publicPath)
	if err != nil {
		return nil, err
	}

	chain := make([]*x509.Certificate, 0)
	if chainPath != "" {
		chain, err = readCertificateChain(ctx, chainPath)
		if err != nil {
			return nil, err
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

	normalizedKey, err := bundle.NormalizePrivateKey(ctx, raw)
	if err != nil {
		return nil, err
	}

	return normalizedKey, nil
}

func readCertificate(ctx context.Context, path string) (*x509.Certificate, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, coreerror.New(
			i18n.M(ctx, i18n.K.CertificateExternalPathNotReadable).V("path", path),
			true,
		)
	}

	parsed, err := bundle.ParseCertificate(ctx, raw)
	if err != nil {
		return nil, err
	}

	return parsed, nil
}

func readCertificateChain(ctx context.Context, path string) ([]*x509.Certificate, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, coreerror.New(
			i18n.M(ctx, i18n.K.CertificateExternalPathNotReadable).V("path", path),
			true,
		)
	}

	certificates, err := bundle.ParseCertificateChain(ctx, raw)
	if err != nil {
		return nil, err
	}

	return certificates, nil
}

func encodeChain(chain []*x509.Certificate) []string {
	result := make([]string, len(chain))
	for index, item := range chain {
		result[index] = base64.StdEncoding.EncodeToString(item.Raw)
	}

	return result
}
