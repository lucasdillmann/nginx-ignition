package letsencrypt

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"

	acmelog "github.com/go-acme/lego/v4/log"

	"dillmann.com.br/nginx-ignition/certificate/commons"
	"dillmann.com.br/nginx-ignition/core/certificate"
	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"dillmann.com.br/nginx-ignition/core/common/coreerror"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

const (
	certificateProviderID = "LETS_ENCRYPT"
	privateKeySize        = 2048
)

type Provider struct {
	configuration *configuration.Configuration
}

func New(cfg *configuration.Configuration) *Provider {
	acmelog.Logger = logAdapterInstance

	return &Provider{
		configuration: cfg,
	}
}

func (p *Provider) ID() string {
	return certificateProviderID
}

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateLetsencryptName)
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

	productionEnvironment, err := p.isProductionEnvironment()
	if err != nil {
		return nil, err
	}

	email, _ := request.Parameters[emailAddressFieldID].(string)

	usrKey, err := rsa.GenerateKey(rand.Reader, privateKeySize)
	if err != nil {
		return nil, coreerror.New(i18n.M(ctx, i18n.K.CertificateLetsencryptGeneratePrivateKey), false)
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
		productionEnvironment,
	)
}

func (p *Provider) Renew(
	ctx context.Context,
	cert *certificate.Certificate,
) (*certificate.Certificate, error) {
	var metadata *certificateMetadata
	if err := json.Unmarshal([]byte(*cert.Metadata), &metadata); err != nil {
		return nil, coreerror.New(i18n.M(ctx, i18n.K.CertificateLetsencryptParseMetadata), false)
	}

	encodedPrivKey, err := base64.StdEncoding.DecodeString(metadata.UserPrivateKey)
	if err != nil {
		return nil, coreerror.New(i18n.M(ctx, i18n.K.CertificateLetsencryptDecodePrivateKey), false)
	}

	privKey, err := x509.ParsePKCS1PrivateKey(encodedPrivKey)
	if err != nil {
		return nil, coreerror.New(i18n.M(ctx, i18n.K.CertificateLetsencryptParsePrivateKey), false)
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
		metadata.ProductionEnvironment,
	)
}

func (p *Provider) isProductionEnvironment() (bool, error) {
	return p.configuration.GetBoolean("nginx-ignition.certificate.lets-encrypt.production")
}
