package acme

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"log/slog"
	"time"

	acmelog "github.com/go-acme/lego/v5/log"

	"github.com/lucasdillmann/nginx-ignition/internal/certificate/commons"
	"github.com/lucasdillmann/nginx-ignition/internal/core/certificate"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/configuration"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/coreerror"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/dynamicfields"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/i18n"
)

const (
	certificateProviderID = "ACME"
	privateKeySize        = 2048
)

type Provider struct {
	configuration *configuration.Configuration
}

func New(cfg *configuration.Configuration) *Provider {
	acmelog.SetDefault(slog.New(&logAdapter{}))

	return &Provider{
		configuration: cfg,
	}
}

func (p *Provider) ID() string {
	return certificateProviderID
}

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateAcmeName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return resolveDynamicFields(ctx)
}

func (p *Provider) Priority() int {
	return 1
}

func (p *Provider) Issue(
	ctx context.Context,
	request *certificate.IssueRequest,
) (*certificate.Certificate, error) {
	if err := commons.Validate(ctx, request, validationRules{p.DynamicFields(ctx)}); err != nil {
		return nil, err
	}

	email, _ := request.Parameters[emailAddressFieldID].(string)

	usrKey, err := rsa.GenerateKey(rand.Reader, privateKeySize)
	if err != nil {
		return nil, coreerror.New(
			i18n.M(ctx, i18n.K.CertificateAcmeGeneratePrivateKey),
			false,
		)
	}

	user := userDetails{
		email:      email,
		privateKey: usrKey,
		newAccount: true,
	}

	return issueCertificate(
		ctx,
		user,
		request.DomainNames,
		request.Parameters,
	)
}

func (p *Provider) Renew(
	ctx context.Context,
	cert *certificate.Certificate,
) (*certificate.Certificate, error) {
	var metadata *certificateMetadata
	if err := json.Unmarshal([]byte(*cert.Metadata), &metadata); err != nil {
		return nil, coreerror.New(i18n.M(ctx, i18n.K.CertificateAcmeParseMetadata), false)
	}

	encodedPrivKey, err := base64.StdEncoding.DecodeString(metadata.UserPrivateKey)
	if err != nil {
		return nil, coreerror.New(i18n.M(ctx, i18n.K.CertificateAcmeDecodePrivateKey), false)
	}

	privKey, err := x509.ParsePKCS1PrivateKey(encodedPrivKey)
	if err != nil {
		return nil, coreerror.New(i18n.M(ctx, i18n.K.CertificateAcmeParsePrivateKey), false)
	}

	user := userDetails{
		email:      metadata.UserMail,
		privateKey: privKey,
		newAccount: false,
	}

	return issueCertificate(
		ctx,
		user,
		cert.DomainNames,
		cert.Parameters,
	)
}

func (p *Provider) IsDueToRenew(
	_ context.Context,
	existing *certificate.Certificate,
) (bool, error) {
	if existing.RenewAfter == nil {
		return false, nil
	}

	return !time.Now().Before(*existing.RenewAfter), nil
}
