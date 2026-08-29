package katapult

import (
	"context"

	"github.com/go-acme/lego/v5/challenge"
	"github.com/go-acme/lego/v5/providers/dns/katapult"

	"nginx-ignition/internal/certificate/letsencrypt/dns"
	"nginx-ignition/internal/core/common/dynamicfields"
	"nginx-ignition/internal/core/common/i18n"
)

const (
	apiKeyFieldID = "katapultApiKey" // nolint:gosec
)

type Provider struct{}

func (p *Provider) ID() string {
	return "KATAPULT"
}

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateLetsencryptDnsKatapultName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          apiKeyFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsKatapultApiKey),
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

	cfg := katapult.NewDefaultConfig()
	cfg.APIKey = apiKey
	cfg.TTL = dns.TTL
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval

	return katapult.NewDNSProviderConfig(cfg)
}
