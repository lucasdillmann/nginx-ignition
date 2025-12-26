package webnames

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/webnames"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

const (
	apiKeyFieldID = "webnamesApiKey"
)

type Provider struct{}

func (p *Provider) ID() string {
	return "WEBNAMES"
}

func (p *Provider) Name() string {
	return "Webnames"
}

func (p *Provider) DynamicFields() []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          apiKeyFieldID,
			Description: "Webnames API key",
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

	cfg := &webnames.Config{
		APIKey:             apiKey,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PollingInterval,
	}

	return webnames.NewDNSProviderConfig(cfg)
}
