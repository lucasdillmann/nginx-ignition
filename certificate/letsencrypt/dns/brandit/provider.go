package brandit

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/brandit"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

//nolint:gosec
const (
	apiKeyFieldID      = "branditAPIKey"
	apiUsernameFieldID = "branditAPIUsername"
)

type Provider struct{}

func (p *Provider) ID() string { return "BRANDIT" }

func (p *Provider) Name() string { return "Brandit" }

func (p *Provider) DynamicFields() []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          apiKeyFieldID,
			Description: "Brandit API key",
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          apiUsernameFieldID,
			Description: "Brandit API username",
			Required:    true,
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
	apiUsername, _ := parameters[apiUsernameFieldID].(string)

	cfg := brandit.NewDefaultConfig()
	cfg.APIKey = apiKey
	cfg.APIUsername = apiUsername
	cfg.TTL = dns.TTL
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval

	return brandit.NewDNSProviderConfig(cfg)
}
