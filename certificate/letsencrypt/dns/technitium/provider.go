package technitium

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/technitium"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

const (
	baseURLFieldID  = "technitiumBaseUrl"
	apiTokenFieldID = "technitiumApiToken"
)

type Provider struct{}

func (p *Provider) ID() string {
	return "TECHNITIUM"
}

func (p *Provider) Name() string {
	return "Technitium"
}

func (p *Provider) DynamicFields() []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          baseURLFieldID,
			Description: "Technitium base URL",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          apiTokenFieldID,
			Description: "Technitium API token",
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
	baseURL, _ := parameters[baseURLFieldID].(string)
	apiToken, _ := parameters[apiTokenFieldID].(string)

	cfg := technitium.NewDefaultConfig()
	cfg.BaseURL = baseURL
	cfg.APIToken = apiToken
	cfg.TTL = dns.TTL
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval

	return technitium.NewDNSProviderConfig(cfg)
}
