package zoneee

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/zoneee"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

const (
	usernameFieldID = "zoneeeUsername"
	apiKeyFieldID   = "zoneeeApiKey"
)

type Provider struct{}

func (p *Provider) ID() string {
	return "ZONEEE"
}

func (p *Provider) Name() string {
	return "Zone.ee"
}

func (p *Provider) DynamicFields() []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          usernameFieldID,
			Description: "Zone.ee username",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          apiKeyFieldID,
			Description: "Zone.ee API key",
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
	username, _ := parameters[usernameFieldID].(string)
	apiKey, _ := parameters[apiKeyFieldID].(string)

	cfg := &zoneee.Config{
		Username:           username,
		APIKey:             apiKey,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PollingInterval,
	}

	return zoneee.NewDNSProviderConfig(cfg)
}
