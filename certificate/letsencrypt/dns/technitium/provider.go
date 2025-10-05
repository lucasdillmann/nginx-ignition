package technitium

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/technitium"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
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

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamic_fields.DynamicField{
		{
			ID:          baseURLFieldID,
			Description: "Technitium base URL",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          apiTokenFieldID,
			Description: "Technitium API token",
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
	baseURL, _ := parameters[baseURLFieldID].(string)
	apiToken, _ := parameters[apiTokenFieldID].(string)

	cfg := &technitium.Config{
		BaseURL:            baseURL,
		APIToken:           apiToken,
		TTL:                dns.TTL,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PollingInterval,
	}

	return technitium.NewDNSProviderConfig(cfg)
}
