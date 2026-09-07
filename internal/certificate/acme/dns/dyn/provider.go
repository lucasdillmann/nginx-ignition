package dyn

import (
	"context"

	"github.com/go-acme/lego/v5/challenge"
	"github.com/go-acme/lego/v5/providers/dns/dyn"

	"github.com/lucasdillmann/nginx-ignition/internal/certificate/acme/dns"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/dynamicfields"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/i18n"
)

const (
	customerNameFieldID = "dynDnsCustomerName"
	usernameFieldID     = "dynDnsUsername"
	passwordFieldID     = "dynDnsPassword"
)

type Provider struct{}

func (p *Provider) ID() string { return "DYN_DNS" }

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateAcmeDnsDynName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          customerNameFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateAcmeDnsDynCustomerName),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          usernameFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateAcmeDnsDynUsername),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          passwordFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateAcmeDnsDynPassword),
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
	customerName, _ := parameters[customerNameFieldID].(string)
	username, _ := parameters[usernameFieldID].(string)
	password, _ := parameters[passwordFieldID].(string)

	cfg := dyn.NewDefaultConfig()
	cfg.CustomerName = customerName
	cfg.UserName = username
	cfg.Password = password
	cfg.TTL = dns.TTL
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval

	return dyn.NewDNSProviderConfig(cfg)
}
