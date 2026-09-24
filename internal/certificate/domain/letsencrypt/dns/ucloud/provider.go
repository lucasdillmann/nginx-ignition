package ucloud

import (
	"context"

	"github.com/go-acme/lego/v5/challenge"
	"github.com/go-acme/lego/v5/providers/dns/ucloud"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/dynamicfields"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/i18n"
	dns2 "github.com/lucasdillmann/nginx-ignition/internal/certificate/domain/letsencrypt/dns"
)

const (
	publicKeyFieldID  = "uCloudPublicKey"
	privateKeyFieldID = "uCloudPrivateKey"
	regionFieldID     = "uCloudRegion"
	projectIDFieldID  = "uCloudProjectId"
)

type Provider struct{}

func (p *Provider) ID() string {
	return "UCLOUD"
}

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateLetsencryptDnsUcloudName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns2.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          publicKeyFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsUcloudPublicKey),
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          privateKeyFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsUcloudPrivateKey),
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          regionFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsUcloudRegion),
			Required:    false,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          projectIDFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsUcloudProjectId),
			Required:    false,
			Type:        dynamicfields.SingleLineTextType,
		},
	})
}

func (p *Provider) ChallengeProvider(
	_ context.Context,
	_ []string,
	parameters map[string]any,
) (challenge.Provider, error) {
	publicKey, _ := parameters[publicKeyFieldID].(string)
	privateKey, _ := parameters[privateKeyFieldID].(string)
	region, _ := parameters[regionFieldID].(string)
	projectID, _ := parameters[projectIDFieldID].(string)

	cfg := ucloud.NewDefaultConfig()
	cfg.PublicKey = publicKey
	cfg.PrivateKey = privateKey
	cfg.Region = region
	cfg.ProjectID = projectID
	cfg.TTL = dns2.TTL
	cfg.PropagationTimeout = dns2.PropagationTimeout
	cfg.PollingInterval = dns2.PollingInterval

	return ucloud.NewDNSProviderConfig(cfg)
}
