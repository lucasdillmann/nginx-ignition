package regfish

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/regfish"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

const (
	apiKeyFieldID = "regfishApiKey"
)

type Provider struct{}

func (p *Provider) ID() string {
	return "REGFISH"
}

func (p *Provider) Name() string {
	return "Regfish"
}

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamic_fields.DynamicField{
		{
			ID:          apiKeyFieldID,
			Description: "Regfish API key",
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

	cfg := &regfish.Config{
		APIKey:             apiKey,
		TTL:                dns.TTL,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PollingInterval,
	}

	return regfish.NewDNSProviderConfig(cfg)
}
