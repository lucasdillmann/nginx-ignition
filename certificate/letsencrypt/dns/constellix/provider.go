package constellix

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/constellix"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

//nolint:gosec
const (
	apiKeyFieldID    = "constellixApiKey"
	secretKeyFieldID = "constellixSecretKey"
)

type Provider struct{}

func (p *Provider) ID() string { return "CONSTELLIX" }

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsConstellixName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          apiKeyFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsConstellixApiKey),
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          secretKeyFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsConstellixSecretKey),
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

	cfg := constellix.NewDefaultConfig()
	cfg.APIKey = apiKey
	cfg.SecretKey = secretKey
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval
	cfg.TTL = dns.TTL

	return constellix.NewDNSProviderConfig(cfg)
}
