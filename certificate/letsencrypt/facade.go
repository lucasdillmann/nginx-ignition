package letsencrypt

import (
	"crypto"
	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
)

type userDetails struct {
	email        string
	registration *registration.Resource
	privateKey   crypto.PrivateKey
}

func (u *userDetails) GetEmail() string {
	return u.email
}

func (u *userDetails) GetRegistration() *registration.Resource {
	return u.registration
}

func (u *userDetails) GetPrivateKey() crypto.PrivateKey {
	return u.privateKey
}

func issueCertificate(
	user userDetails,
	domainNames []string,
	parameters map[string]interface{},
	configuration configuration.Configuration,
) (*certificate.Resource, error) {
	caURL := "https://acme-v02.api.letsencrypt.org/directory"

	production, err := configuration.GetBoolean("nginx-ignition.certificate.lets-encrypt.production")
	if err != nil {
		return nil, err
	}

	if !production {
		caURL = "https://acme-staging-v02.api.letsencrypt.org/directory"
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

	err = client.Challenge.SetDNS01Provider(dnsChallengeProvider)
	if err != nil {
		return nil, err
	}

	if user.registration == nil {
		user.registration, err = client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
		if err != nil {
			return nil, err
		}
	}

	request := certificate.ObtainRequest{
		Domains: domainNames,
		Bundle:  true,
	}

	return client.Certificate.Obtain(request)
}
