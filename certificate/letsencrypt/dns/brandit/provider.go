package brandit

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/brandit"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

const (
	apiKeyFieldID      = "branditAPIKey"
	apiUsernameFieldID = "branditAPIUsername"
)

type Provider struct{}

func (p *Provider) ID() string { return "BRANDIT" }

func (p *Provider) Name() string { return "Brandit" }

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamic_fields.DynamicField{
		{
			ID:          apiKeyFieldID,
			Description: "Brandit API key",
			Required:    true,
			Sensitive:   true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          apiUsernameFieldID,
			Description: "Brandit API username",
			Required:    true,
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
	apiUsername, _ := parameters[apiUsernameFieldID].(string)

	cfg := &brandit.Config{
		APIKey:             apiKey,
		APIUsername:        apiUsername,
		TTL:                dns.TTL,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PollingInterval,
	}

	return brandit.NewDNSProviderConfig(cfg)
}
