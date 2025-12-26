package constellix

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/constellix"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

//nolint:gosec
const (
	apiKeyFieldID    = "constellixApiKey"
	secretKeyFieldID = "constellixSecretKey"
)

type Provider struct{}

func (p *Provider) ID() string { return "CONSTELLIX" }

func (p *Provider) Name() string { return "Constellix" }

func (p *Provider) DynamicFields() []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          apiKeyFieldID,
			Description: "Constellix API key",
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          secretKeyFieldID,
			Description: "Constellix secret key",
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

	cfg := &constellix.Config{
		APIKey:             apiKey,
		SecretKey:          secretKey,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PollingInterval,
		TTL:                dns.TTL,
	}

	return constellix.NewDNSProviderConfig(cfg)
}
