package active24

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/active24"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

const (
	apiKeyFieldID = "active24ApiKey"
	secretFieldID = "active24Secret"
)

type Provider struct{}

func (p *Provider) ID() string { return "ACTIVE24" }

func (p *Provider) Name() string { return "Active24" }

func (p *Provider) DynamicFields() []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          apiKeyFieldID,
			Description: "Active24 API key",
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          secretFieldID,
			Description: "Active24 secret",
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
	secret, _ := parameters[secretFieldID].(string)

	cfg := &active24.Config{
		APIKey:             apiKey,
		Secret:             secret,
		TTL:                dns.TTL,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PollingInterval,
	}

	return active24.NewDNSProviderConfig(cfg)
}
