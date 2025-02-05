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
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
	"time"
)

type Provider struct {
	configuration configuration.Configuration
}

func New(configuration configuration.Configuration) *Provider {
	return &Provider{
		configuration: configuration,
	}
}

func (p *Provider) ID() string {
	return "LETS_ENCRYPT"
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
	productionEnvironment, err := p.configuration.GetBoolean("nginx-ignition.certificate.lets-encrypt.production")
	if err != nil {
		return nil, err
	}

	email, casted := request.Parameters[emailAddress.ID].(string)
	if !casted {
		return nil, core_error.New("E-mail address is missing", true)
	}

	usrKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, core_error.New("Failed to generate private key", false)
	}

	user := userDetails{
		email:      email,
		privateKey: usrKey,
	}

	cert, err := issueCertificate(
		user,
		request.DomainNames,
		request.Parameters,
		p.configuration,
	)

	if err != nil {
		return nil, err
	}

	certDetails, err := x509.ParseCertificate(cert.Certificate)
	if err != nil {
		return nil, core_error.New("Failed to parse CSR", false)
	}

	metadata := certificateMetadata{
		UserMail:              email,
		UserPrivateKey:        base64.StdEncoding.EncodeToString(x509.MarshalPKCS1PrivateKey(usrKey)),
		UserPublicKey:         base64.StdEncoding.EncodeToString(x509.MarshalPKCS1PublicKey(&usrKey.PublicKey)),
		ProductionEnvironment: productionEnvironment,
	}
	metadataJson, err := jsoniter.MarshalToString(metadata)
	if err != nil {
		return nil, core_error.New("Failed to serialize metadata", false)
	}

	renewAfter := time.Now().Add(5 * time.Hour * 24)
	output := certificate.Certificate{
		ID:                 uuid.New(),
		ProviderID:         p.ID(),
		DomainNames:        request.DomainNames,
		IssuedAt:           time.Now(),
		ValidUntil:         certDetails.NotAfter,
		ValidFrom:          certDetails.NotBefore,
		RenewAfter:         &renewAfter,
		PrivateKey:         base64.StdEncoding.EncodeToString(cert.PrivateKey),
		PublicKey:          base64.StdEncoding.EncodeToString(certDetails.RawSubjectPublicKeyInfo),
		CertificationChain: []string{base64.StdEncoding.EncodeToString(cert.IssuerCertificate)},
		Parameters:         request.Parameters,
		Metadata:           &metadataJson,
	}

	return &output, nil
}

func (p *Provider) Renew(_ *certificate.Certificate) (*certificate.Certificate, error) {
	// TODO: Implement this
	return nil, core_error.New("not implemented yet", false)
}
