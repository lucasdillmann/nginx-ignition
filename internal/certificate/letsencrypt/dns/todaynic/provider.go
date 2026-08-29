package todaynic

import (
	"context"

	"github.com/go-acme/lego/v5/challenge"
	"github.com/go-acme/lego/v5/providers/dns/todaynic"

	"nginx-ignition/internal/certificate/letsencrypt/dns"
	"nginx-ignition/internal/core/common/dynamicfields"
	"nginx-ignition/internal/core/common/i18n"
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
