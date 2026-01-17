package tencentcloud

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/tencentcloud"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
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
	return i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsTencentcloudName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          secretIDKeyFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsTencentcloudSecretId),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          secretKeyKeyFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsTencentcloudSecretKey),
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          regionFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsTencentcloudRegion),
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
	cfg.TTL = dns.TTL
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval

	return tencentcloud.NewDNSProviderConfig(cfg)
}
