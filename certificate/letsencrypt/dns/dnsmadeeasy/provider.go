package dnsmadeeasy

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/dnsmadeeasy"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

//nolint:gosec
const (
	apiKeyFieldID    = "dnsMadeEasyApiKey"
	secretKeyFieldID = "dnsMadeEasySecretKey"
)

type Provider struct{}

func (p *Provider) ID() string { return "DNSMADEEASY" }

func (p *Provider) Name() string { return "DNSMadeEasy" }

func (p *Provider) DynamicFields() []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          apiKeyFieldID,
			Description: "DNSMadeEasy API key",
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          secretKeyFieldID,
			Description: "DNSMadeEasy secret key",
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
	secretKey, _ := parameters[secretKeyFieldID].(string)

	cfg := &dnsmadeeasy.Config{
		APIKey:             apiKey,
		APISecret:          secretKey,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PollingInterval,
		TTL:                dns.TTL,
	}

	return dnsmadeeasy.NewDNSProviderConfig(cfg)
}
