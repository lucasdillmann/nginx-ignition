package conohav3

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/conohav3"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

const (
	userIDFieldID   = "conohaV3UserId"
	passwordFieldID = "conohaV3Password"
	tenantIDFieldID = "conohaV3TenantId"
	regionFieldID   = "conohaV3Region"
)

type Provider struct{}

func (p *Provider) ID() string { return "CONOHA_V3" }

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateLetsencryptDnsConohav3Name)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          userIDFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsConohav3UserId),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          passwordFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsConohav3Password),
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          tenantIDFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsConohav3TenantId),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          regionFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsConohav3Region),
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
	userName, _ := parameters[userIDFieldID].(string)
	password, _ := parameters[passwordFieldID].(string)
	tenantID, _ := parameters[tenantIDFieldID].(string)
	region, _ := parameters[regionFieldID].(string)

	cfg := conohav3.NewDefaultConfig()
	cfg.UserID = userName
	cfg.Password = password
	cfg.TenantID = tenantID
	cfg.Region = region
	cfg.TTL = dns.TTL
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval

	return conohav3.NewDNSProviderConfig(cfg)
}
