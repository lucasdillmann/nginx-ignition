package vegadns

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/vegadns"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

const (
	baseURLFieldID   = "vegadnsBaseUrl"
	apiKeyFieldID    = "vegadnsApiKey"
	apiSecretFieldID = "vegadnsApiSecret"
)

type Provider struct{}

func (p *Provider) ID() string {
	return "VEGADNS"
}

func (p *Provider) Name() string {
	return "VegaDNS"
}

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamic_fields.DynamicField{
		{
			ID:          baseURLFieldID,
			Description: "VegaDNS base URL",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          apiKeyFieldID,
			Description: "VegaDNS API key",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          apiSecretFieldID,
			Description: "VegaDNS API secret",
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
	apiKey, _ := parameters[apiKeyFieldID].(string)
	apiSecret, _ := parameters[apiSecretFieldID].(string)

	cfg := &vegadns.Config{
		BaseURL:            baseURL,
		APIKey:             apiKey,
		APISecret:          apiSecret,
		TTL:                dns.TTL,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PoolingInterval,
	}

	return vegadns.NewDNSProviderConfig(cfg)
}
