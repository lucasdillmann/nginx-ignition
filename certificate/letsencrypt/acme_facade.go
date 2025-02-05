package letsencrypt

import (
	"crypto/x509"
	"dillmann.com.br/nginx-ignition/core/certificate"
	"dillmann.com.br/nginx-ignition/core/common/core_error"
	"encoding/base64"
	"encoding/pem"
	"github.com/go-acme/lego/v4/certcrypto"
	acmecertificate "github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
	"strings"
	"time"
)

func issueCertificate(
	user userDetails,
	domainNames []string,
	parameters map[string]any,
	productionEnvironment bool,
) (*certificate.Certificate, error) {
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

	dnsChallengeProvider, err := resolveDnsProvider(parameters)
	if err != nil {
		return nil, err
	}

	client.Challenge.Remove(challenge.TLSALPN01)
	client.Challenge.Remove(challenge.HTTP01)
	err = client.Challenge.SetDNS01Provider(dnsChallengeProvider)
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

	return parseResult(uuid.New(), domainNames, parameters, cert, user, productionEnvironment, client)
}

func parseResult(
	id uuid.UUID,
	domainNames []string,
	parameters map[string]any,
	result *acmecertificate.Resource,
	usr userDetails,
	productionEnvironment bool,
	client *lego.Client,
) (*certificate.Certificate, error) {
	mainCert, _ := strings.CutSuffix(string(result.Certificate), string(result.IssuerCertificate))
	mainCertBytes := []byte(mainCert)

	pemBlock, _ := pem.Decode(mainCertBytes)
	if pemBlock == nil || pemBlock.Type != "CERTIFICATE" {
		return nil, core_error.New("failed to decode PEM block containing certificate", false)
	}

	certDetails, err := x509.ParseCertificate(pemBlock.Bytes)
	if err != nil {
		return nil, core_error.New("Failed to parse CSR", false)
	}

	metadata := certificateMetadata{
		UserMail:              usr.email,
		UserPrivateKey:        base64.StdEncoding.EncodeToString(x509.MarshalPKCS1PrivateKey(usr.privateKey)),
		UserPublicKey:         base64.StdEncoding.EncodeToString(x509.MarshalPKCS1PublicKey(&usr.privateKey.PublicKey)),
		ProductionEnvironment: productionEnvironment,
	}
	metadataJson, err := jsoniter.MarshalToString(metadata)
	if err != nil {
		return nil, core_error.New("Failed to serialize metadata", false)
	}

	renewalInfo, err := client.Certificate.GetRenewalInfo(acmecertificate.RenewalInfoRequest{certDetails})
	if err != nil {
		return nil, core_error.New("Failed to get renewal info", false)
	}

	output := certificate.Certificate{
		ID:                 id,
		ProviderID:         certificateProviderId,
		DomainNames:        domainNames,
		IssuedAt:           time.Now(),
		ValidUntil:         certDetails.NotAfter,
		ValidFrom:          certDetails.NotBefore,
		RenewAfter:         &renewalInfo.SuggestedWindow.Start,
		PrivateKey:         base64.StdEncoding.EncodeToString(result.PrivateKey),
		PublicKey:          base64.StdEncoding.EncodeToString(mainCertBytes),
		CertificationChain: []string{base64.StdEncoding.EncodeToString(result.IssuerCertificate)},
		Parameters:         parameters,
		Metadata:           &metadataJson,
	}

	return &output, nil
}
