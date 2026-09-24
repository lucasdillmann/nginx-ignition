package conohav2

import (
	"context"

	"github.com/go-acme/lego/v5/challenge"
	"github.com/go-acme/lego/v5/providers/dns/conoha"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/dynamicfields"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/i18n"
	dns2 "github.com/lucasdillmann/nginx-ignition/internal/certificate/domain/letsencrypt/dns"
)

const (
	regionFieldID   = "conoHaV2Region"
	tenantIDFieldID = "conoHaV2TenantID"
	usernameFieldID = "conoHaV2Username"
	passwordFieldID = "conoHaV2Password"
)

type Provider struct{}

func (p *Provider) ID() string { return "CONOHA_V2" }

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateLetsencryptDnsConohav2Name)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns2.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          regionFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsConohav2Region),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          tenantIDFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsConohav2TenantId),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          usernameFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsConohav2Username),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          passwordFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsConohav2Password),
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
	})
}

func (p *Provider) ChallengeProvider(
	_ context.Context,
	_ []string,
	parameters map[string]any,
) (challenge.Provider, error) {
	region, _ := parameters[regionFieldID].(string)
	tenantID, _ := parameters[tenantIDFieldID].(string)
	username, _ := parameters[usernameFieldID].(string)
	password, _ := parameters[passwordFieldID].(string)

	cfg := conoha.NewDefaultConfig()
	cfg.Region = region
	cfg.TenantID = tenantID
	cfg.Username = username
	cfg.Password = password
	cfg.PropagationTimeout = dns2.PropagationTimeout
	cfg.PollingInterval = dns2.PollingInterval
	cfg.TTL = dns2.TTL

	return conoha.NewDNSProviderConfig(cfg)
}
