package clouddns

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/clouddns"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

const (
	clientIDFieldID = "cloudDnsClientId"
	emailFieldID    = "cloudDnsEmail"
	passwordFieldID = "cloudDnsPassword"
)

type Provider struct{}

func (p *Provider) ID() string { return "CLOUDDNS" }

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsClouddnsName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          clientIDFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsClouddnsClientId),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          emailFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsClouddnsEmail),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          passwordFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsClouddnsPassword),
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
	clientID, _ := parameters[clientIDFieldID].(string)
	email, _ := parameters[emailFieldID].(string)
	password, _ := parameters[passwordFieldID].(string)

	cfg := clouddns.NewDefaultConfig()
	cfg.ClientID = clientID
	cfg.Email = email
	cfg.Password = password
	cfg.TTL = dns.TTL
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval

	return clouddns.NewDNSProviderConfig(cfg)
}
