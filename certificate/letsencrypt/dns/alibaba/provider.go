package alibaba

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/alidns"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
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
	return i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsAlibabaName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          accessKeyFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsAlibabaAccessKeyId),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          accessKeySecretFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsAlibabaAccessKeySecret),
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          securityTokenFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsAlibabaSecurityToken),
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          regionFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsAlibabaRegion),
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          ramRoleFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsAlibabaRamRole),
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
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval
	cfg.TTL = dns.TTL

	return alidns.NewDNSProviderConfig(cfg)
}
