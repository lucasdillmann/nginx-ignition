package tencentcloud

import (
	"context"

	"github.com/go-acme/lego/v5/challenge"
	"github.com/go-acme/lego/v5/providers/dns/tencentcloud"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/dynamicfields"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/i18n"
	dns2 "github.com/lucasdillmann/nginx-ignition/internal/certificate/domain/letsencrypt/dns"
)

//nolint:gosec
const (
	secretIDKeyFieldID  = "tencentCloudSecretId"
	secretKeyKeyFieldID = "tencentCloudSecretKey"
	regionFieldID       = "tencentCloudRegion"
)

type Provider struct{}

func (p *Provider) ID() string { return "TENCENT_CLOUD" }

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateLetsencryptDnsTencentcloudName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns2.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          secretIDKeyFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsTencentcloudSecretId),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          secretKeyKeyFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsTencentcloudSecretKey),
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          regionFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsTencentcloudRegion),
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
	secretID, _ := parameters[secretIDKeyFieldID].(string)
	secretKey, _ := parameters[secretKeyKeyFieldID].(string)
	region, _ := parameters[regionFieldID].(string)

	cfg := tencentcloud.NewDefaultConfig()
	cfg.SecretID = secretID
	cfg.SecretKey = secretKey
	cfg.Region = region
	cfg.TTL = dns2.TTL
	cfg.PropagationTimeout = dns2.PropagationTimeout
	cfg.PollingInterval = dns2.PollingInterval

	return tencentcloud.NewDNSProviderConfig(cfg)
}
