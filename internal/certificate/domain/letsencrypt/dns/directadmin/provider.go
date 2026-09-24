package directadmin

import (
	"context"

	"github.com/go-acme/lego/v5/challenge"
	"github.com/go-acme/lego/v5/providers/dns/directadmin"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/dynamicfields"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/i18n"
	dns2 "github.com/lucasdillmann/nginx-ignition/internal/certificate/domain/letsencrypt/dns"
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
	return dns2.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
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
	cfg.TTL = dns2.TTL
	cfg.PropagationTimeout = dns2.PropagationTimeout
	cfg.PollingInterval = dns2.PollingInterval
	cfg.HTTPClient = nil

	return directadmin.NewDNSProviderConfig(cfg)
}
