package bundle

import (
	"bytes"
	"context"
	"crypto"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/coreerror"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/i18n"
)

func ParsePrivateKey(ctx context.Context, raw []byte) (crypto.Signer, error) {
	if len(bytes.TrimSpace(raw)) == 0 {
		return nil, invalidPrivateKey(ctx)
	}

	foundPEM := false
	remaining := raw

	for {
		block, rest := pem.Decode(remaining)
		if block == nil {
			break
		}

		foundPEM = true
		remaining = rest

		key, isPrivateKey, err := parsePrivateKeyBlock(ctx, block)
		if err != nil {
			return nil, err
		}

		if isPrivateKey {
			return key, nil
		}
	}

	if foundPEM {
		return nil, invalidPrivateKey(ctx)
	}

	key, isPrivateKey, err := parsePrivateKeyDER(ctx, raw)
	if err != nil {
		return nil, err
	}

	if !isPrivateKey {
		return nil, invalidPrivateKey(ctx)
	}

	return key, nil
}

func NormalizePrivateKey(ctx context.Context, raw []byte) ([]byte, error) {
	key, err := ParsePrivateKey(ctx, raw)
	if err != nil {
		return nil, err
	}

	normalized, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return nil, coreerror.New(
			i18n.M(ctx, i18n.K.CertificateBundleUnsupportedPrivateKey),
			true,
		)
	}

	return normalized, nil
}

func ParseCertificate(ctx context.Context, raw []byte) (*x509.Certificate, error) {
	certificates, err := ParseCertificateChain(ctx, raw)
	if err != nil {
		return nil, err
	}

	return certificates[0], nil
}

func ParseCertificateChain(ctx context.Context, raw []byte) ([]*x509.Certificate, error) {
	if len(bytes.TrimSpace(raw)) == 0 {
		return nil, invalidCertificate(ctx)
	}

	certificates := make([]*x509.Certificate, 0)
	foundPEM := false
	remaining := raw

	for {
		block, rest := pem.Decode(remaining)
		if block == nil {
			break
		}

		foundPEM = true
		remaining = rest

		if !isCertificateBlock(block.Type) {
			continue
		}

		certificate, err := parseCertificateBlock(ctx, block)
		if err != nil {
			return nil, err
		}

		certificates = append(certificates, certificate)
	}

	if len(certificates) > 0 {
		return certificates, nil
	}

	if foundPEM {
		return nil, invalidCertificate(ctx)
	}

	parsed, err := x509.ParseCertificates(raw)
	if err != nil || len(parsed) == 0 {
		return nil, invalidCertificate(ctx)
	}

	return parsed, nil
}

func parsePrivateKeyBlock(ctx context.Context, block *pem.Block) (crypto.Signer, bool, error) {
	if isEncryptedPrivateKeyBlock(block) {
		return nil, true, coreerror.New(
			i18n.M(ctx, i18n.K.CertificateBundleEncryptedPrivateKey),
			true,
		)
	}

	switch block.Type {
	case "DSA PRIVATE KEY", "OPENSSH PRIVATE KEY":
		return nil, true, coreerror.New(
			i18n.M(ctx, i18n.K.CertificateBundleUnsupportedPrivateKey),
			true,
		)
	}

	return parsePrivateKeyDER(ctx, block.Bytes)
}

func parsePrivateKeyDER(ctx context.Context, der []byte) (crypto.Signer, bool, error) {
	pkcs8Key, pkcs8Err := x509.ParsePKCS8PrivateKey(der)
	if pkcs8Err == nil {
		signer, ok := pkcs8Key.(crypto.Signer)
		if !ok {
			return nil, true, coreerror.New(
				i18n.M(ctx, i18n.K.CertificateBundleUnsupportedPrivateKey),
				true,
			)
		}

		return signer, true, nil
	}

	pkcs1Key, pkcs1Err := x509.ParsePKCS1PrivateKey(der)
	if pkcs1Err == nil {
		return pkcs1Key, true, nil
	}

	sec1Key, sec1Err := x509.ParseECPrivateKey(der)
	if sec1Err == nil {
		return sec1Key, true, nil
	}

	return nil, false, nil
}

func isEncryptedPrivateKeyBlock(block *pem.Block) bool {
	if block.Type == "ENCRYPTED PRIVATE KEY" {
		return true
	}

	return block.Headers["Proc-Type"] == "4,ENCRYPTED"
}

func isCertificateBlock(blockType string) bool {
	switch blockType {
	case "CERTIFICATE", "X509 CERTIFICATE", "TRUSTED CERTIFICATE":
		return true
	default:
		return false
	}
}

func parseCertificateBlock(ctx context.Context, block *pem.Block) (*x509.Certificate, error) {
	der := block.Bytes
	if block.Type == "TRUSTED CERTIFICATE" {
		der = trimTrustedCertificate(der)
	}

	certificate, err := x509.ParseCertificate(der)
	if err != nil {
		return nil, invalidCertificate(ctx)
	}

	return certificate, nil
}

func trimTrustedCertificate(der []byte) []byte {
	var value asn1.RawValue

	rest, err := asn1.Unmarshal(der, &value)
	if err != nil || len(rest) == 0 {
		return der
	}

	return value.FullBytes
}

func invalidPrivateKey(ctx context.Context) *coreerror.CoreError {
	msg := i18n.M(ctx, i18n.K.CertificateBundleInvalidPrivateKey)
	return coreerror.New(msg, true)
}

func invalidCertificate(ctx context.Context) *coreerror.CoreError {
	msg := i18n.M(ctx, i18n.K.CertificateBundleInvalidCertificate)
	return coreerror.New(msg, true)
}
