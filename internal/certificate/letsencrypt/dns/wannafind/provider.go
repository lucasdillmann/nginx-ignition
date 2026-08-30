package wannafind

import (
	"context"

	"github.com/go-acme/lego/v5/challenge"
	"github.com/go-acme/lego/v5/providers/dns/wannafind"

	"github.com/lucasdillmann/nginx-ignition/internal/certificate/letsencrypt/dns"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/dynamicfields"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/i18n"
)

const (
	apiKeyFieldID = "wannafindApiKey"
)

type Provider struct{}

func (p *Provider) ID() string {
	return "WANNAFIND"
}

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateLetsencryptDnsWannafindName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          apiKeyFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsWannafindApiKey),
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

	cfg := wannafind.NewDefaultConfig()
	cfg.APIKey = apiKey
	cfg.TTL = dns.TTL
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval

	return wannafind.NewDNSProviderConfig(cfg)
}
