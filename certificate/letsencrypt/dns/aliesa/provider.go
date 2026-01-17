package aliesa

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/aliesa"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

const (
	ramRoleFieldID       = "aliesaRamRole"
	apiKeyFieldID        = "aliesaApiKey"
	secretKeyFieldID     = "aliesaSecretKey"
	securityTokenFieldID = "aliesaSecurityToken"
	regionIDFieldID      = "aliesaRegionID"
)

type Provider struct{}

func (p *Provider) ID() string {
	return "ALIESA"
}

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsAliesaName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          ramRoleFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsAliesaRamRole),
			Required:    false,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          apiKeyFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsAliesaApiKey),
			Required:    false,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          secretKeyFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsAliesaSecretKey),
			Required:    false,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          securityTokenFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsAliesaSecurityToken),
			Required:    false,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          regionIDFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsAliesaRegionId),
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
	ramRole, _ := parameters[ramRoleFieldID].(string)
	apiKey, _ := parameters[apiKeyFieldID].(string)
	secretKey, _ := parameters[secretKeyFieldID].(string)
	securityToken, _ := parameters[securityTokenFieldID].(string)
	regionID, _ := parameters[regionIDFieldID].(string)

	cfg := aliesa.NewDefaultConfig()
	cfg.RAMRole = ramRole
	cfg.APIKey = apiKey
	cfg.SecretKey = secretKey
	cfg.SecurityToken = securityToken
	cfg.RegionID = regionID
	cfg.TTL = dns.TTL
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval

	return aliesa.NewDNSProviderConfig(cfg)
}
