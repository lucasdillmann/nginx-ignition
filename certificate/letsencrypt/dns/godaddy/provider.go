package godaddy

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/godaddy"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

//nolint:gosec
const (
	apiKeyFieldID    = "goDaddyApiKey"
	apiSecretFieldID = "goDaddyApiSecret"
)

type Provider struct{}

func (p *Provider) ID() string { return "GODADDY" }

func (p *Provider) Name() string { return "GoDaddy" }

func (p *Provider) DynamicFields() []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          apiKeyFieldID,
			Description: "GoDaddy API key",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          apiSecretFieldID,
			Description: "GoDaddy API secret",
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
	key, _ := parameters[apiKeyFieldID].(string)
	secret, _ := parameters[apiSecretFieldID].(string)

	cfg := &godaddy.Config{
		APIKey:             key,
		APISecret:          secret,
		TTL:                dns.TTL,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PollingInterval,
	}

	return godaddy.NewDNSProviderConfig(cfg)
}
