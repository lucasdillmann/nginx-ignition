package ucloud

import (
	"context"

	"github.com/go-acme/lego/v5/challenge"
	"github.com/go-acme/lego/v5/providers/dns/ucloud"

	"github.com/lucasdillmann/nginx-ignition/internal/certificate/acme/dns"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/dynamicfields"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/i18n"
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
	return i18n.M(ctx, i18n.K.CertificateAcmeDnsUcloudName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          publicKeyFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateAcmeDnsUcloudPublicKey),
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          privateKeyFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateAcmeDnsUcloudPrivateKey),
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          regionFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateAcmeDnsUcloudRegion),
			Required:    false,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          projectIDFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateAcmeDnsUcloudProjectId),
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
	cfg.TTL = dns.TTL
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval

	return ucloud.NewDNSProviderConfig(cfg)
}
