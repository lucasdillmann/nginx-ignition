package vegadns

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/vegadns"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

//nolint:gosec
const (
	baseURLFieldID   = "vegadnsBaseUrl"
	apiKeyFieldID    = "vegadnsApiKey"
	apiSecretFieldID = "vegadnsApiSecret"
)

type Provider struct{}

func (p *Provider) ID() string {
	return "VEGADNS"
}

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsVegadnsName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          baseURLFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsVegadnsBaseUrl),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          apiKeyFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsVegadnsApiKey),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          apiSecretFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsVegadnsApiSecret),
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
	apiSecret, _ := parameters[apiSecretFieldID].(string)

	cfg := vegadns.NewDefaultConfig()
	cfg.BaseURL = baseURL
	cfg.APIKey = apiKey
	cfg.APISecret = apiSecret
	cfg.TTL = dns.TTL
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval

	return vegadns.NewDNSProviderConfig(cfg)
}
