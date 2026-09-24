package nicru

import (
	"context"

	"github.com/go-acme/lego/v5/challenge"
	"github.com/go-acme/lego/v5/providers/dns/nicru"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/dynamicfields"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/i18n"
	dns2 "github.com/lucasdillmann/nginx-ignition/internal/certificate/domain/letsencrypt/dns"
)

//nolint:gosec
const (
	usernameFieldID  = "nicRuUsername"
	passwordFieldID  = "nicRuPassword"
	serviceIDFieldID = "nicRuServiceId"
	secretFieldID    = "nicRuSecret"
)

type Provider struct{}

func (p *Provider) ID() string { return "NICRU" }

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateLetsencryptDnsNicruName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns2.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          usernameFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsNicruUsername),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          passwordFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsNicruPassword),
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          serviceIDFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsNicruOauth2ServiceId),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          secretFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsNicruOauth2Secret),
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
	username, _ := parameters[usernameFieldID].(string)
	password, _ := parameters[passwordFieldID].(string)
	serviceID, _ := parameters[serviceIDFieldID].(string)
	secret, _ := parameters[secretFieldID].(string)

	cfg := nicru.NewDefaultConfig()
	cfg.Username = username
	cfg.Password = password
	cfg.ServiceID = serviceID
	cfg.Secret = secret
	cfg.TTL = dns2.TTL
	cfg.PropagationTimeout = dns2.PropagationTimeout
	cfg.PollingInterval = dns2.PollingInterval

	return nicru.NewDNSProviderConfig(cfg)
}
