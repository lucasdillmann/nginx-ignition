package porkbun

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/porkbun"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

const (
	apiKeyFieldID       = "porkbunApiKey"
	secretApiKeyFieldID = "porkbunSecretApiKey"
)

type Provider struct{}

func (p *Provider) ID() string {
	return "PORKBUN"
}

func (p *Provider) Name() string {
	return "Porkbun"
}

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamic_fields.DynamicField{
		{
			ID:          apiKeyFieldID,
			Description: "Porkbun API key",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          secretApiKeyFieldID,
			Description: "Porkbun secret API key",
			Required:    true,
			Sensitive:   true,
			Type:        dynamic_fields.SingleLineTextType,
		},
	})
}

func (p *Provider) ChallengeProvider(
	_ context.Context,
	_ []string,
	parameters map[string]any,
) (challenge.Provider, error) {
	apiKey, _ := parameters[apiKeyFieldID].(string)
	secretApiKey, _ := parameters[secretApiKeyFieldID].(string)

	cfg := &porkbun.Config{
		APIKey:             apiKey,
		SecretAPIKey:       secretApiKey,
		TTL:                dns.TTL,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PollingInterval,
	}

	return porkbun.NewDNSProviderConfig(cfg)
}
