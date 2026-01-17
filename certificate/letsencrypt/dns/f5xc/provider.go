package f5xc

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/f5xc"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

//nolint:gosec
const (
	apiTokenFieldID   = "f5xcApiToken"
	tenantNameFieldID = "f5xcTenantName"
	groupNameFieldID  = "f5xcGroupName"
)

type Provider struct{}

func (p *Provider) ID() string { return "F5_XC" }

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsF5xcName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          apiTokenFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsF5xcApiToken),
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          tenantNameFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsF5xcTenantName),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          groupNameFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsF5xcGroupName),
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
	apiToken, _ := parameters[apiTokenFieldID].(string)
	tenantName, _ := parameters[tenantNameFieldID].(string)
	groupName, _ := parameters[groupNameFieldID].(string)

	cfg := f5xc.NewDefaultConfig()
	cfg.APIToken = apiToken
	cfg.TenantName = tenantName
	cfg.GroupName = groupName
	cfg.TTL = dns.TTL
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval

	return f5xc.NewDNSProviderConfig(cfg)
}
