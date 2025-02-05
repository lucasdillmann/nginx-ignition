package letsencrypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"dillmann.com.br/nginx-ignition/core/certificate"
	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"dillmann.com.br/nginx-ignition/core/common/core_error"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
	"encoding/base64"
	"encoding/json"
	acmelog "github.com/go-acme/lego/v4/log"
)

const (
	certificateProviderId = "LETS_ENCRYPT"
	privateKeySize        = 2048
)

type Provider struct {
	configuration *configuration.Configuration
}

func New(configuration *configuration.Configuration) *Provider {
	acmelog.Logger = logAdapterInstance

	return &Provider{
		configuration: configuration,
	}
}

func (p *Provider) ID() string {
	return certificateProviderId
}

func (p *Provider) Name() string {
	return "Let's encrypt"
}

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return []*dynamic_fields.DynamicField{
		&termsOfService,
		&emailAddress,
		&dnsProvider,
		&awsAccessKey,
		&awsSecretKey,
		&awsHostedZoneID,
		&cloudflareApiToken,
		&googleCloudPrivateKey,
		&azureTenantId,
		&azureSubscriptionId,
		&azureClientId,
		&azureClientSecret,
		&azureEnvironment,
	}
}

func (p *Provider) Priority() int {
	return 1
}

func (p *Provider) Issue(request *certificate.IssueRequest) (*certificate.Certificate, error) {
	productionEnvironment, err := p.isProductionEnvironment()
	if err != nil {
		return nil, err
	}

	email, casted := request.Parameters[emailAddress.ID].(string)
	if !casted {
		return nil, core_error.New("E-mail address is missing", true)
	}

	usrKey, err := rsa.GenerateKey(rand.Reader, privateKeySize)
	if err != nil {
		return nil, core_error.New("Failed to generate private key", false)
	}

	user := userDetails{
		email:      email,
		privateKey: usrKey,
		newAccount: true,
	}

	return issueCertificate(
		user,
		request.DomainNames,
		request.Parameters,
		productionEnvironment,
	)
}

func (p *Provider) Renew(cert *certificate.Certificate) (*certificate.Certificate, error) {
	var metadata *certificateMetadata
	if err := json.Unmarshal([]byte(*cert.Metadata), &metadata); err != nil {
		return nil, core_error.New("Failed to parse metadata", false)
	}

	encodedPrivKey, err := base64.StdEncoding.DecodeString(metadata.UserPrivateKey)
	if err != nil {
		return nil, core_error.New("Failed to decode private key", false)
	}

	privKey, err := x509.ParsePKCS1PrivateKey(encodedPrivKey)
	if err != nil {
		return nil, core_error.New("Failed to parse private key", false)
	}

	user := userDetails{
		email:      metadata.UserMail,
		privateKey: privKey,
		newAccount: false,
	}

	return issueCertificate(user, cert.DomainNames, cert.Parameters, metadata.ProductionEnvironment)
}

func (p *Provider) isProductionEnvironment() (bool, error) {
	return p.configuration.GetBoolean("nginx-ignition.certificate.lets-encrypt.production")
}
