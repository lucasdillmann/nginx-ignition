package baiducloud

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/baiducloud"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

//nolint:gosec
const (
	accessKeyFieldID       = "baiduCloudAccessKeyID"
	secretAccessKeyFieldID = "baiduCloudSecretAccessKey"
)

type Provider struct{}

func (p *Provider) ID() string { return "BAIDU_CLOUD" }

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateLetsencryptDnsBaiducloudName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          accessKeyFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsBaiducloudAccessKeyId),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID: secretAccessKeyFieldID,
			Description: i18n.M(
				ctx,
				i18n.K.CertificateLetsencryptDnsBaiducloudSecretAccessKey,
			),
			Required:  true,
			Sensitive: true,
			Type:      dynamicfields.SingleLineTextType,
		},
	})
}

func (p *Provider) ChallengeProvider(
	_ context.Context,
	_ []string,
	parameters map[string]any,
) (challenge.Provider, error) {
	accessKey, _ := parameters[accessKeyFieldID].(string)
	secretKey, _ := parameters[secretAccessKeyFieldID].(string)

	cfg := baiducloud.NewDefaultConfig()
	cfg.AccessKeyID = accessKey
	cfg.SecretAccessKey = secretKey
	cfg.TTL = dns.TTL
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval

	return baiducloud.NewDNSProviderConfig(cfg)
}
