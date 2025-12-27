package cloudflare

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/cloudflare"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

const (
	apiTokenFieldID = "cloudflareApiToken"
)

type Provider struct{}

func (p *Provider) ID() string {
	return "CLOUDFLARE"
}

func (p *Provider) Name() string {
	return "Cloudflare"
}

func (p *Provider) DynamicFields() []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          apiTokenFieldID,
			Description: "Cloudflare API token",
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

	cfg := &cloudflare.Config{
		AuthToken:          apiToken,
		TTL:                dns.TTL,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PollingInterval,
	}

	return cloudflare.NewDNSProviderConfig(cfg)
}
