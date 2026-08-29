package letsencrypt

import (
	"context"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"time"

	"github.com/go-acme/lego/v5/certcrypto"
	acmecertificate "github.com/go-acme/lego/v5/certificate"
	"github.com/go-acme/lego/v5/challenge"
	"github.com/go-acme/lego/v5/challenge/dns01"
	"github.com/go-acme/lego/v5/lego"
	"github.com/go-acme/lego/v5/registration"
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"

	"nginx-ignition/internal/core/certificate"
	"nginx-ignition/internal/core/common/coreerror"
	"nginx-ignition/internal/core/common/i18n"
)

func issueCertificate(
	ctx context.Context,
	user userDetails,
	domainNames []string,
	parameters map[string]any,
	productionEnvironment bool,
) (*certificate.Certificate, error) {
	caURL := lego.DirectoryURLLetsEncrypt
	if !productionEnvironment {
		caURL = lego.DirectoryURLLetsEncryptStaging
	}

	config := lego.NewConfig(&user)
	config.CADirURL = caURL

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

	dnsOptions := make([]dns01.ChallengeOption, 0, 1)
	propagationBypassCasted, propagationBypass := parameters["bypassDnsPropagationChecks"].(bool)
	if propagationBypassCasted && propagationBypass {
		dnsOptions = append(dnsOptions, dns01.PropagationWait(5*time.Second, true))
	}

	err = client.Challenge.SetDNS01Provider(dnsChallenge, dnsOptions...)
	if err != nil {
		return nil, err
	}

	registerOptions := registration.RegisterOptions{TermsOfServiceAgreed: true}
	if user.newAccount {
		user.registration, err = client.Registration.Register(ctx, registerOptions)
	} else {
		user.registration, err = client.Registration.ResolveAccountByKey(ctx)
		if err != nil {
			user.registration, err = client.Registration.Register(ctx, registerOptions)
		}
	}

	if err != nil {
		return nil, err
	}

	request := acmecertificate.ObtainRequest{
		Domains: domainNames,
		Bundle:  true,
		KeyType: certcrypto.RSA2048,
	}
	if productionEnvironment {
		request.PreferredChain = "ISRG Root X1"
	}

	cert, err := client.Certificate.Obtain(ctx, request)
	if err != nil {
		return nil, err
	}

	return parseResult(
		ctx,
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
	ctx context.Context,
	id uuid.UUID,
	domainNames []string,
	parameters map[string]any,
	result *acmecertificate.Resource,
	usr userDetails,
	productionEnvironment bool,
	client *lego.Client,
) (*certificate.Certificate, error) {
	pemBlock, _ := pem.Decode(result.Certificate)
	if pemBlock == nil || pemBlock.Type != "CERTIFICATE" {
		return nil, coreerror.New(
			i18n.M(ctx, i18n.K.CommonUnableToParsePem).V("type", "certificate"),
			false,
		)
	}

	metadata := certificateMetadata{
		UserMail: usr.email,
		UserPrivateKey: base64.StdEncoding.EncodeToString(
			x509.MarshalPKCS1PrivateKey(usr.privateKey),
		),
		UserPublicKey: base64.StdEncoding.EncodeToString(
			x509.MarshalPKCS1PublicKey(&usr.privateKey.PublicKey),
		),
		ProductionEnvironment: productionEnvironment,
	}

	metadataJSON, err := jsoniter.MarshalToString(metadata)
	if err != nil {
		return nil, err
	}

	privateKey, err := certcrypto.ParsePEMPrivateKey(result.PrivateKey)
	if err != nil {
		return nil, coreerror.New(
			i18n.M(ctx, i18n.K.CommonUnableToParsePem).V("type", "private key"),
			false,
		)
	}

	encodedPrivateKey, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return nil, err
	}

	notAfter, notBefore, renewAt, err := fetchCertDates(ctx, *pemBlock, client)
	if err != nil {
		return nil, err
	}

	certificationChain, err := encodeIssuerCertificates(ctx, result.IssuerCertificate)
	if err != nil {
		return nil, err
	}

	output := certificate.Certificate{
		ID:                 id,
		ProviderID:         certificateProviderID,
		DomainNames:        domainNames,
		IssuedAt:           time.Now(),
		ValidUntil:         *notAfter,
		ValidFrom:          *notBefore,
		RenewAfter:         renewAt,
		PrivateKey:         base64.StdEncoding.EncodeToString(encodedPrivateKey),
		PublicKey:          base64.StdEncoding.EncodeToString(pemBlock.Bytes),
		CertificationChain: certificationChain,
		Parameters:         parameters,
		Metadata:           &metadataJSON,
	}

	return &output, nil
}

func encodeIssuerCertificates(ctx context.Context, issuer []byte) ([]string, error) {
	chain := make([]string, 0)
	remaining := issuer

	for {
		pemBlock, rest := pem.Decode(remaining)
		if pemBlock == nil {
			break
		}

		remaining = rest

		if pemBlock.Type != "CERTIFICATE" {
			continue
		}

		chain = append(chain, base64.StdEncoding.EncodeToString(pemBlock.Bytes))
	}

	if len(chain) == 0 {
		return nil, coreerror.New(
			i18n.M(ctx, i18n.K.CommonUnableToParsePem).V("type", "issuer"),
			false,
		)
	}

	return chain, nil
}

func fetchCertDates(
	ctx context.Context,
	pemBlock pem.Block,
	client *lego.Client,
) (
	notAfter *time.Time,
	notBefore *time.Time,
	renewAt *time.Time,
	err error,
) {
	certDetails, err := x509.ParseCertificate(pemBlock.Bytes)
	if err != nil {
		return nil, nil, nil, err
	}

	renewalInfo, err := client.Certificate.GetRenewalInfo(ctx, certDetails)
	if err != nil {
		return nil, nil, nil, err
	}

	notAfter = &certDetails.NotAfter
	notBefore = &certDetails.NotBefore
	renewAt = &renewalInfo.SuggestedWindow.Start

	return notAfter, notBefore, renewAt, nil
}
