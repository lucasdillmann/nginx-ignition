package domeneshop

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/domeneshop"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

const (
	apiTokenFieldID  = "domeneshopApiToken"
	apiSecretFieldID = "domeneshopApiSecret"
)

type Provider struct{}

func (p *Provider) ID() string { return "DOMENESHOP" }

func (p *Provider) Name() string { return "Domeneshop" }

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamic_fields.DynamicField{
		{
			ID:          apiTokenFieldID,
			Description: "Domeneshop API token",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          apiSecretFieldID,
			Description: "Domeneshop API secret",
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
	apiToken, _ := parameters[apiTokenFieldID].(string)
	apiSecret, _ := parameters[apiSecretFieldID].(string)

	cfg := &domeneshop.Config{
		APIToken:           apiToken,
		APISecret:          apiSecret,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PoolingInterval,
	}

	return domeneshop.NewDNSProviderConfig(cfg)
}
