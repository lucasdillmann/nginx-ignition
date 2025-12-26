package exoscale

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/exoscale"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

const (
	apiKeyFieldID    = "exoscaleAPIKey"
	apiSecretFieldID = "exoscaleAPISecret"
)

type Provider struct{}

func (p *Provider) ID() string { return "EXOSCALE" }

func (p *Provider) Name() string { return "Exoscale" }

func (p *Provider) DynamicFields() []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          apiKeyFieldID,
			Description: "Exoscale API key",
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          apiSecretFieldID,
			Description: "Exoscale API secret",
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
	apiSecret, _ := parameters[apiSecretFieldID].(string)

	cfg := &exoscale.Config{
		APIKey:             apiKey,
		APISecret:          apiSecret,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PollingInterval,
		TTL:                dns.TTL,
	}

	return exoscale.NewDNSProviderConfig(cfg)
}
