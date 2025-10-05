package auroradns

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/auroradns"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

const (
	apiKeyFieldID = "auroraDnsApiKey"
	secretFieldID = "auroraDnsSecret"
)

type Provider struct{}

func (p *Provider) ID() string { return "AURORA_DNS" }

func (p *Provider) Name() string { return "Aurora DNS" }

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamic_fields.DynamicField{
		{
			ID:          apiKeyFieldID,
			Description: "Aurora DNS API key",
			Required:    true,
			Sensitive:   true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          secretFieldID,
			Description: "Aurora DNS secret",
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
	apiKey, _ := parameters[apiKeyFieldID].(string)
	secret, _ := parameters[secretFieldID].(string)

	cfg := &auroradns.Config{
		APIKey:             apiKey,
		Secret:             secret,
		TTL:                dns.TTL,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PoolingInterval,
	}

	return auroradns.NewDNSProviderConfig(cfg)
}
