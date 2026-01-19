package sakuracloud

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/sakuracloud"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

//nolint:gosec
const (
	accessTokenFieldID  = "sakuraCloudAccessToken"
	accessSecretFieldID = "sakuraCloudAccessSecret"
)

type Provider struct{}

func (p *Provider) ID() string { return "SAKURA_CLOUD" }

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateLetsencryptDnsSakuracloudName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          accessTokenFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsSakuracloudAccessToken),
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          accessSecretFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsSakuracloudAccessSecret),
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
	accessToken, _ := parameters[accessTokenFieldID].(string)
	accessSecret, _ := parameters[accessSecretFieldID].(string)

	cfg := sakuracloud.NewDefaultConfig()
	cfg.Token = accessToken
	cfg.Secret = accessSecret
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval
	cfg.TTL = dns.TTL

	return sakuracloud.NewDNSProviderConfig(cfg)
}
