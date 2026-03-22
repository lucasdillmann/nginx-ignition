package namesurfer

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/namesurfer"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

const (
	apiKeyFieldID    = "nameSurferApiKey" // nolint:gosec
	apiSecretFieldID = "nameSurferApiSecret"
	baseURLFieldID   = "nameSurferBaseUrl"
)

type Provider struct{}

func (p *Provider) ID() string {
	return "NAMESURFER"
}

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateLetsencryptDnsNamesurferName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          apiKeyFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsNamesurferApiKey),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          apiSecretFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsNamesurferApiSecret),
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          baseURLFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsNamesurferBaseUrl),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
	})
}

func (p *Provider) ChallengeProvider(
	_ context.Context,
	_ []string,
	parameters map[string]any,
) (challenge.Provider, error) {
	config := namesurfer.NewDefaultConfig()
	config.APIKey, _ = parameters[apiKeyFieldID].(string)
	config.APISecret, _ = parameters[apiSecretFieldID].(string)
	config.BaseURL, _ = parameters[baseURLFieldID].(string)

	return namesurfer.NewDNSProviderConfig(config)
}
