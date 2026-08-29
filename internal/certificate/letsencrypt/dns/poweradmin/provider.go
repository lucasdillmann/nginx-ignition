package poweradmin

import (
	"context"

	"github.com/go-acme/lego/v5/challenge"
	"github.com/go-acme/lego/v5/providers/dns/poweradmin"

	"nginx-ignition/internal/certificate/letsencrypt/dns"
	"nginx-ignition/internal/core/common/dynamicfields"
	"nginx-ignition/internal/core/common/i18n"
)

//nolint:gosec
const (
	baseURLFieldID = "poweradminBaseUrl"
	apiKeyFieldID  = "poweradminApiKey"
)

type Provider struct{}

func (p *Provider) ID() string {
	return "POWERADMIN"
}

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateLetsencryptDnsPoweradminName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          baseURLFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsPoweradminBaseUrl),
			Required:    true,
			Type:        dynamicfields.URLType,
		},
		{
			ID:          apiKeyFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsPoweradminApiKey),
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
	baseURL, _ := parameters[baseURLFieldID].(string)
	apiKey, _ := parameters[apiKeyFieldID].(string)

	cfg := poweradmin.NewDefaultConfig()
	cfg.BaseURL = baseURL
	cfg.APIKey = apiKey
	cfg.TTL = dns.TTL
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval

	return poweradmin.NewDNSProviderConfig(cfg)
}
