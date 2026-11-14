package cloudpress

import (
	"context"

	"github.com/go-acme/lego/v5/challenge"
	"github.com/go-acme/lego/v5/providers/dns/cloudpress"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/dynamicfields"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/i18n"
	"github.com/lucasdillmann/nginx-ignition/internal/certificate/domain/letsencrypt/dns"
)

//nolint:gosec
const (
	baseURLFieldID   = "cloudPressBaseURL"
	apiKeyFieldID    = "cloudPressApiKey"
	accountIDFieldID = "cloudPressAccountID"
)

type Provider struct{}

func (p *Provider) ID() string {
	return "CLOUDPRESS"
}

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateLetsencryptDnsCloudpressName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          baseURLFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsCloudpressBaseUrl),
			Required:    true,
			Type:        dynamicfields.URLType,
		},
		{
			ID:          apiKeyFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsCloudpressApiKey),
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          accountIDFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsCloudpressAccountId),
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
	accountID, _ := parameters[accountIDFieldID].(string)

	cfg := cloudpress.NewDefaultConfig()
	cfg.BaseURL = baseURL
	cfg.APIKey = apiKey
	cfg.AccountID = accountID
	cfg.TTL = dns.TTL
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval

	return cloudpress.NewDNSProviderConfig(cfg)
}
