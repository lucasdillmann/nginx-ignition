package aliesa

import (
	"context"

	"github.com/go-acme/lego/v5/challenge"
	"github.com/go-acme/lego/v5/providers/dns/aliesa"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/dynamicfields"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/i18n"
	dns2 "github.com/lucasdillmann/nginx-ignition/internal/certificate/domain/letsencrypt/dns"
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
	return i18n.M(ctx, i18n.K.CertificateLetsencryptDnsAliesaName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns2.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          ramRoleFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsAliesaRamRole),
			Required:    false,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          apiKeyFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsAliesaApiKey),
			Required:    false,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          secretKeyFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsAliesaSecretKey),
			Required:    false,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          securityTokenFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsAliesaSecurityToken),
			Required:    false,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          regionIDFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsAliesaRegionId),
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
	cfg.TTL = dns2.TTL
	cfg.PropagationTimeout = dns2.PropagationTimeout
	cfg.PollingInterval = dns2.PollingInterval

	return aliesa.NewDNSProviderConfig(cfg)
}
