package f5xc

import (
	"context"

	"github.com/go-acme/lego/v5/challenge"
	"github.com/go-acme/lego/v5/providers/dns/f5xc"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/dynamicfields"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/i18n"
	dns2 "github.com/lucasdillmann/nginx-ignition/internal/certificate/domain/letsencrypt/dns"
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
	return i18n.M(ctx, i18n.K.CertificateLetsencryptDnsF5xcName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns2.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          apiTokenFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsF5xcApiToken),
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          tenantNameFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsF5xcTenantName),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          groupNameFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsF5xcGroupName),
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
	cfg.TTL = dns2.TTL
	cfg.PropagationTimeout = dns2.PropagationTimeout
	cfg.PollingInterval = dns2.PollingInterval

	return f5xc.NewDNSProviderConfig(cfg)
}
