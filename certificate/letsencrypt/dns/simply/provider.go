package simply

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/simply"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

const (
	accountNameFieldID = "simplyAccountName"
	apiKeyFieldID      = "simplyApiKey"
)

type Provider struct{}

func (p *Provider) ID() string {
	return "SIMPLY"
}

func (p *Provider) Name() string {
	return "Simply.com"
}

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamic_fields.DynamicField{
		{
			ID:          accountNameFieldID,
			Description: "Simply.com account name",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          apiKeyFieldID,
			Description: "Simply.com API key",
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
	accountName, _ := parameters[accountNameFieldID].(string)
	apiKey, _ := parameters[apiKeyFieldID].(string)

	cfg := &simply.Config{
		AccountName:        accountName,
		APIKey:             apiKey,
		TTL:                dns.TTL,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PollingInterval,
	}

	return simply.NewDNSProviderConfig(cfg)
}
