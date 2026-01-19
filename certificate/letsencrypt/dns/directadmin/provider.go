package directadmin

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/directadmin"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

const (
	hostFieldID     = "directAdminHost"
	passwordFieldID = "directAdminPassword"
	userFieldID     = "directAdminUsername"
	zoneNameFieldID = "directAdminZoneName"
)

type Provider struct{}

func (p *Provider) ID() string { return "DIRECTADMIN" }

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateLetsencryptDnsDirectadminName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          hostFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsDirectadminBaseUrl),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          userFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsDirectadminUsername),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          passwordFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsDirectadminPassword),
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          zoneNameFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsDirectadminZoneName),
			Type:        dynamicfields.SingleLineTextType,
		},
	})
}

func (p *Provider) ChallengeProvider(
	_ context.Context,
	_ []string,
	parameters map[string]any,
) (challenge.Provider, error) {
	host, _ := parameters[hostFieldID].(string)
	user, _ := parameters[userFieldID].(string)
	password, _ := parameters[passwordFieldID].(string)
	zoneName, _ := parameters[zoneNameFieldID].(string)

	cfg := directadmin.NewDefaultConfig()
	cfg.BaseURL = host
	cfg.Username = user
	cfg.Password = password
	cfg.ZoneName = zoneName
	cfg.TTL = dns.TTL
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval
	cfg.HTTPClient = nil

	return directadmin.NewDNSProviderConfig(cfg)
}
