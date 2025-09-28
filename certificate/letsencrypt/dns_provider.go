package letsencrypt

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns/aws"
	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns/azure"
	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns/cloudflare"
	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns/gcp"
	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns/porkbun"
	"dillmann.com.br/nginx-ignition/core/common/core_error"
)

var (
	providers = []dns.Provider{
		&aws.Provider{},
		&azure.Provider{},
		&cloudflare.Provider{},
		&gcp.Provider{},
		&porkbun.Provider{},
	}
)

func resolveProviderChallenge(ctx context.Context, domainNames []string, parameters map[string]any) (challenge.Provider, error) {
	providerId, _ := parameters[dnsProvider.ID].(string)

	for _, provider := range providers {
		if provider.ID() == providerId {
			return provider.ChallengeProvider(ctx, domainNames, parameters)
		}
	}

	return nil, core_error.New("Unknown DNS provider", true)
}
