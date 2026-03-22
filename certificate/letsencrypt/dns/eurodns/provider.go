package eurodns

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/eurodns"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

const (
	applicationIDFieldID = "euroDnsApplicationId"
	apiKeyFieldID        = "euroDnsApiKey" // nolint:gosec
)

type Provider struct{}

func (p *Provider) ID() string {
	return "EURO_DNS"
}

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateLetsencryptDnsEurodnsName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          applicationIDFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsEurodnsApplicationId),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          apiKeyFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsEurodnsApiKey),
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
	config := eurodns.NewDefaultConfig()
	config.ApplicationID, _ = parameters[applicationIDFieldID].(string)
	config.APIKey, _ = parameters[apiKeyFieldID].(string)

	return eurodns.NewDNSProviderConfig(config)
}
