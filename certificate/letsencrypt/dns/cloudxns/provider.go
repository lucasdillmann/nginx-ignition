package cloudxns

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/cloudxns"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

//nolint:gosec
const (
	apiKeyFieldID    = "cloudXnsApiKey"
	secretKeyFieldID = "cloudXnsSecretKey"
)

type Provider struct{}

func (p *Provider) ID() string { return "CLOUDXNS" }

func (p *Provider) Name() string { return "CloudXNS" }

func (p *Provider) DynamicFields() []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          apiKeyFieldID,
			Description: "CloudXNS API key",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          secretKeyFieldID,
			Description: "CloudXNS secret key",
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

	cfg := &cloudxns.Config{
		APIKey:             apiKey,
		SecretKey:          secretKey,
		TTL:                dns.TTL,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PollingInterval,
	}

	return cloudxns.NewDNSProviderConfig(cfg)
}
