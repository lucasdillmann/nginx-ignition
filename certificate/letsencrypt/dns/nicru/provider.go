package nicru

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/nicru"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
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
	return i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsNicruName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          usernameFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsNicruUsername),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          passwordFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsNicruPassword),
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          serviceIDFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsNicruOauth2ServiceId),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          secretFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsNicruOauth2Secret),
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
	cfg.TTL = dns.TTL
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval

	return nicru.NewDNSProviderConfig(cfg)
}
