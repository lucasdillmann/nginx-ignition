package alwaysdata

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/alwaysdata"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

const (
	apiKeyFieldID  = "alwaysdataApiKey"
	accountFieldID = "alwaysdataAccount"
)

type Provider struct{}

func (p *Provider) ID() string {
	return "ALWAYSDATA"
}

func (p *Provider) Name() string {
	return "Alwaysdata"
}

func (p *Provider) DynamicFields() []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          apiKeyFieldID,
			Description: "Alwaysdata API key",
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          accountFieldID,
			Description: "Alwaysdata account name (optional)",
			Required:    false,
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
	account, _ := parameters[accountFieldID].(string)

	cfg := alwaysdata.NewDefaultConfig()
	cfg.APIKey = apiKey
	cfg.Account = account
	cfg.TTL = dns.TTL
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval

	return alwaysdata.NewDNSProviderConfig(cfg)
}
