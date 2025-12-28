package porkbun

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/porkbun"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

//nolint:gosec
const (
	apiKeyFieldID       = "porkbunApiKey"
	secretAPIKeyFieldID = "porkbunSecretApiKey"
)

type Provider struct{}

func (p *Provider) ID() string {
	return "PORKBUN"
}

func (p *Provider) Name() string {
	return "Porkbun"
}

func (p *Provider) DynamicFields() []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          apiKeyFieldID,
			Description: "Porkbun API key",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          secretAPIKeyFieldID,
			Description: "Porkbun secret API key",
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
	})
}

func (p *Provider) ChallengeProvider(
	_ context.Context,
	_ []string,
	parameters map[string]any,
) (challenge.Provider, error) {
	apiKey, _ := parameters[apiKeyFieldID].(string)
	secretAPIKey, _ := parameters[secretAPIKeyFieldID].(string)

	cfg := &porkbun.Config{
		APIKey:             apiKey,
		SecretAPIKey:       secretAPIKey,
		TTL:                dns.TTL,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PollingInterval,
	}

	return porkbun.NewDNSProviderConfig(cfg)
}
