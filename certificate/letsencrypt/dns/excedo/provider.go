package excedo

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/excedo"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

const (
	apiKeyFieldID = "excedoApiKey" // nolint:gosec
	apiURLFieldID = "excedoApiUrl"
)

type Provider struct{}

func (p *Provider) ID() string {
	return "EXCEDO"
}

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateLetsencryptDnsExcedoName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          apiKeyFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsExcedoApiKey),
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          apiURLFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsExcedoApiUrl),
			Required:    false,
			Type:        dynamicfields.SingleLineTextType,
		},
	})
}

func (p *Provider) ChallengeProvider(
	_ context.Context,
	_ []string,
	parameters map[string]any,
) (challenge.Provider, error) {
	config := excedo.NewDefaultConfig()
	config.APIKey, _ = parameters[apiKeyFieldID].(string)
	if apiURL, ok := parameters[apiURLFieldID].(string); ok && apiURL != "" {
		config.APIURL = apiURL
	}

	return excedo.NewDNSProviderConfig(config)
}
