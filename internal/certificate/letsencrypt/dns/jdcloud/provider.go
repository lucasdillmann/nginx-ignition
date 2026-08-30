package jdcloud

import (
	"context"

	"github.com/go-acme/lego/v5/challenge"
	"github.com/go-acme/lego/v5/providers/dns/jdcloud"

	"github.com/lucasdillmann/nginx-ignition/internal/certificate/letsencrypt/dns"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/dynamicfields"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/i18n"
)

const (
	accessKeyFieldID = "jdcloudAccessKey"
	secretKeyFieldID = "jdcloudSecretKey"
	regionIDFieldID  = "jdcloudRegionId"
)

type Provider struct{}

func (p *Provider) ID() string {
	return "JD_CLOUD"
}

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateLetsencryptDnsJdcloudName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          accessKeyFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsJdcloudAccessKey),
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          secretKeyFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsJdcloudSecretKey),
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          regionIDFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsJdcloudRegionId),
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
	secretKey, _ := parameters[secretKeyFieldID].(string)
	regionID, _ := parameters[regionIDFieldID].(string)

	cfg := jdcloud.NewDefaultConfig()
	cfg.AccessKeyID = accessKey
	cfg.AccessKeySecret = secretKey
	cfg.RegionID = regionID
	cfg.TTL = dns.TTL
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval

	return jdcloud.NewDNSProviderConfig(cfg)
}
