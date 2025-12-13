package letsencrypt

import (
	"context"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"strings"
	"time"

	"github.com/go-acme/lego/v4/certcrypto"
	acmecertificate "github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"

	"dillmann.com.br/nginx-ignition/core/certificate/server"
	"dillmann.com.br/nginx-ignition/core/common/coreerror"
)

func issueCertificate(
	ctx context.Context,
	user userDetails,
	domainNames []string,
	parameters map[string]any,
	productionEnvironment bool,
) (*server.Certificate, error) {
	caURL := lego.LEDirectoryProduction
	if !productionEnvironment {
		caURL = lego.LEDirectoryStaging
	}

	config := lego.NewConfig(&user)
	config.CADirURL = caURL
	config.Certificate.KeyType = certcrypto.RSA2048

	client, err := lego.NewClient(config)
	if err != nil {
		return nil, err
	}

	dnsChallenge, err := resolveProviderChallenge(ctx, domainNames, parameters)
	if err != nil {
		return nil, err
	}

	client.Challenge.Remove(challenge.TLSALPN01)
	client.Challenge.Remove(challenge.HTTP01)
	err = client.Challenge.SetDNS01Provider(dnsChallenge)
	if err != nil {
		return nil, err
	}

	registerOptions := registration.RegisterOptions{TermsOfServiceAgreed: true}
	if user.newAccount {
		user.registration, err = client.Registration.Register(registerOptions)
	} else {
		user.registration, err = client.Registration.ResolveAccountByKey()
		if err != nil {
			user.registration, err = client.Registration.Register(registerOptions)
		}
	}

	if err != nil {
		return nil, err
	}

	request := acmecertificate.ObtainRequest{
		Domains: domainNames,
		Bundle:  true,
	}

	cert, err := client.Certificate.Obtain(request)
	if err != nil {
		return nil, err
	}

	return parseResult(
		uuid.New(),
		domainNames,
		parameters,
		cert,
		user,
		productionEnvironment,
		client,
	)
}

func parseResult(
	id uuid.UUID,
	domainNames []string,
	parameters map[string]any,
	result *acmecertificate.Resource,
	usr userDetails,
	productionEnvironment bool,
	client *lego.Client,
) (*server.Certificate, error) {
	mainCert := strings.Replace(string(result.Certificate), string(result.IssuerCertificate), "", 1)
	pemBlock, _ := pem.Decode([]byte(mainCert))
	if pemBlock == nil || pemBlock.Type != "CERTIFICATE" {
		return nil, coreerror.New("failed to decode PEM block containing certificate", false)
	}

	metadata := certificateMetadata{
		UserMail:              usr.email,
		UserPrivateKey:        base64.StdEncoding.EncodeToString(x509.MarshalPKCS1PrivateKey(usr.privateKey)),
		UserPublicKey:         base64.StdEncoding.EncodeToString(x509.MarshalPKCS1PublicKey(&usr.privateKey.PublicKey)),
		ProductionEnvironment: productionEnvironment,
	}

	metadataJson, err := jsoniter.MarshalToString(metadata)
	if err != nil {
		return nil, err
	}

	privateKeyBlock, _ := pem.Decode(result.PrivateKey)
	if privateKeyBlock == nil || privateKeyBlock.Type != "RSA PRIVATE KEY" {
		return nil, coreerror.New("failed to decode PEM block with the private key", false)
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		return nil, err
	}

	encodedPrivateKey, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return nil, err
	}

	notAfter, notBefore, renewAt, err := fetchCertDates(*pemBlock, client)
	if err != nil {
		return nil, err
	}

	encodedCertificationChain, err := encodeIssuerCertificate(result.IssuerCertificate)
	if err != nil {
		return nil, err
	}

	output := server.Certificate{
		ID:                 id,
		ProviderID:         certificateProviderId,
		DomainNames:        domainNames,
		IssuedAt:           time.Now(),
		ValidUntil:         *notAfter,
		ValidFrom:          *notBefore,
		RenewAfter:         renewAt,
		PrivateKey:         base64.StdEncoding.EncodeToString(encodedPrivateKey),
		PublicKey:          base64.StdEncoding.EncodeToString(pemBlock.Bytes),
		CertificationChain: []string{*encodedCertificationChain},
		Parameters:         parameters,
		Metadata:           &metadataJson,
	}

	return &output, nil
}

func encodeIssuerCertificate(issuer []byte) (*string, error) {
	pemBlock, _ := pem.Decode(issuer)
	if pemBlock == nil || pemBlock.Type != "CERTIFICATE" {
		return nil, coreerror.New("Failed to decode issuer PEM block", false)
	}

	encodedValue := base64.StdEncoding.EncodeToString(pemBlock.Bytes)
	return &encodedValue, nil
}

func fetchCertDates(pemBlock pem.Block, client *lego.Client) (
	notAfter *time.Time,
	notBefore *time.Time,
	renewAt *time.Time,
	err error,
) {
	certDetails, err := x509.ParseCertificate(pemBlock.Bytes)
	if err != nil {
		return nil, nil, nil, err
	}

	infoRequest := acmecertificate.RenewalInfoRequest{Cert: certDetails}
	renewalInfo, err := client.Certificate.GetRenewalInfo(infoRequest)
	if err != nil {
		return nil, nil, nil, err
	}

	notAfter = &certDetails.NotAfter
	notBefore = &certDetails.NotBefore
	renewAt = &renewalInfo.SuggestedWindow.Start

	return notAfter, notBefore, renewAt, nil
}
