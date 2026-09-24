package huaweicloud

import (
	"context"

	"github.com/go-acme/lego/v5/challenge"
	"github.com/go-acme/lego/v5/providers/dns/huaweicloud"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/dynamicfields"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/i18n"
	dns2 "github.com/lucasdillmann/nginx-ignition/internal/certificate/domain/letsencrypt/dns"
)

//nolint:gosec
const (
	accessKeyFieldID       = "huaweiCloudAccessKeyID"
	secretAccessKeyFieldID = "huaweiCloudSecretAccessKey"
	regionFieldID          = "huaweiCloudRegion"
)

type Provider struct{}

func (p *Provider) ID() string { return "HUAWEI_CLOUD" }

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateLetsencryptDnsHuaweicloudName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns2.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          accessKeyFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsHuaweicloudAccessKeyId),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID: secretAccessKeyFieldID,
			Description: i18n.M(
				ctx,
				i18n.K.CertificateLetsencryptDnsHuaweicloudSecretAccessKey,
			),
			Required:  true,
			Sensitive: true,
			Type:      dynamicfields.SingleLineTextType,
		},
		{
			ID:          regionFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsHuaweicloudRegion),
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
	accessKey, _ := parameters[accessKeyFieldID].(string)
	secretAccessKey, _ := parameters[secretAccessKeyFieldID].(string)
	region, _ := parameters[regionFieldID].(string)

	cfg := huaweicloud.NewDefaultConfig()
	cfg.AccessKeyID = accessKey
	cfg.SecretAccessKey = secretAccessKey
	cfg.Region = region
	cfg.PropagationTimeout = dns2.PropagationTimeout
	cfg.PollingInterval = dns2.PollingInterval
	cfg.TTL = dns2.TTL

	return huaweicloud.NewDNSProviderConfig(cfg)
}
