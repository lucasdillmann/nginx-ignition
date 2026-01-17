package conohav2

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/conoha"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
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
	return i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsConohav2Name)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          regionFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsConohav2Region),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          tenantIDFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsConohav2TenantId),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          usernameFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsConohav2Username),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          passwordFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsConohav2Password),
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
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval
	cfg.TTL = dns.TTL

	return conoha.NewDNSProviderConfig(cfg)
}
