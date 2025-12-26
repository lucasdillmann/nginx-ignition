package websupport

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/websupport"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

const (
	apiKeyFieldID = "websupportApiKey"
	secretFieldID = "websupportSecret"
)

type Provider struct{}

func (p *Provider) ID() string { return "WEBSUPPORT" }

func (p *Provider) Name() string { return "Websupport" }

func (p *Provider) DynamicFields() []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          apiKeyFieldID,
			Description: "Websupport API key",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          secretFieldID,
			Description: "Websupport secret",
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

	cfg := &websupport.Config{
		APIKey:             apiKey,
		Secret:             secret,
		TTL:                dns.TTL,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PollingInterval,
	}

	return websupport.NewDNSProviderConfig(cfg)
}
