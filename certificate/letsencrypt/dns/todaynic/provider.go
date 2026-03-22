package todaynic

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/todaynic"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

const (
	authUserIDFieldID = "todaynicAuthUserId"
	apiKeyFieldID     = "todaynicApiKey" // nolint:gosec
)

type Provider struct{}

func (p *Provider) ID() string {
	return "TODAYNIC"
}

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateLetsencryptDnsTodaynicName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          authUserIDFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsTodaynicAuthUserId),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          apiKeyFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsTodaynicApiKey),
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
	config := todaynic.NewDefaultConfig()
	config.AuthUserID, _ = parameters[authUserIDFieldID].(string)
	config.APIKey, _ = parameters[apiKeyFieldID].(string)

	return todaynic.NewDNSProviderConfig(config)
}
