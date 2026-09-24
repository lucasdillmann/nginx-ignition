package alibaba

import (
	"context"

	"github.com/go-acme/lego/v5/challenge"
	"github.com/go-acme/lego/v5/providers/dns/alidns"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/dynamicfields"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/i18n"
	dns2 "github.com/lucasdillmann/nginx-ignition/internal/certificate/domain/letsencrypt/dns"
)

//nolint:gosec
const (
	accessKeyFieldID       = "alibabaAccessKeyId"
	accessKeySecretFieldID = "alibabaAccessKeySecret"
	securityTokenFieldID   = "alibabaSecurityToken"
	regionFieldID          = "alibabaRegion"
	ramRoleFieldID         = "alibabaRamRole"
)

type Provider struct{}

func (p *Provider) ID() string { return "ALIBABA" }

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateLetsencryptDnsAlibabaName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns2.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          accessKeyFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsAlibabaAccessKeyId),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          accessKeySecretFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsAlibabaAccessKeySecret),
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          securityTokenFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsAlibabaSecurityToken),
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          regionFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsAlibabaRegion),
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          ramRoleFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsAlibabaRamRole),
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
	accessSecret, _ := parameters[accessKeySecretFieldID].(string)
	securityToken, _ := parameters[securityTokenFieldID].(string)
	region, _ := parameters[regionFieldID].(string)
	role, _ := parameters[ramRoleFieldID].(string)

	cfg := alidns.NewDefaultConfig()
	cfg.RAMRole = role
	cfg.APIKey = accessKey
	cfg.SecretKey = accessSecret
	cfg.SecurityToken = securityToken
	cfg.RegionID = region
	cfg.PropagationTimeout = dns2.PropagationTimeout
	cfg.PollingInterval = dns2.PollingInterval
	cfg.TTL = dns2.TTL

	return alidns.NewDNSProviderConfig(cfg)
}
