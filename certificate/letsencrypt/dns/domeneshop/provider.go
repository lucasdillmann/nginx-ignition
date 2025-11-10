package domeneshop

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/domeneshop"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

const (
	apiTokenFieldID  = "domeneshopApiToken"
	apiSecretFieldID = "domeneshopApiSecret"
)

type Provider struct{}

func (p *Provider) ID() string { return "DOMENESHOP" }

func (p *Provider) Name() string { return "Domeneshop" }

func (p *Provider) DynamicFields() []*dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          apiTokenFieldID,
			Description: "Domeneshop API token",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          apiSecretFieldID,
			Description: "Domeneshop API secret",
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
	apiToken, _ := parameters[apiTokenFieldID].(string)
	apiSecret, _ := parameters[apiSecretFieldID].(string)

	cfg := &domeneshop.Config{
		APIToken:           apiToken,
		APISecret:          apiSecret,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PollingInterval,
	}

	return domeneshop.NewDNSProviderConfig(cfg)
}
