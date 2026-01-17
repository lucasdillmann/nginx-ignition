package spaceship

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/spaceship"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

//nolint:gosec
const (
	apiKeyFieldID    = "spaceshipApiKey"
	apiSecretFieldID = "spaceshipApiSecret"
)

type Provider struct{}

func (p *Provider) ID() string {
	return "SPACESHIP"
}

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsSpaceshipName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          apiKeyFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsSpaceshipApiKey),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          apiSecretFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsSpaceshipApiSecret),
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
	apiSecret, _ := parameters[apiSecretFieldID].(string)

	cfg := spaceship.NewDefaultConfig()
	cfg.APIKey = apiKey
	cfg.APISecret = apiSecret
	cfg.TTL = dns.TTL
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval

	return spaceship.NewDNSProviderConfig(cfg)
}
