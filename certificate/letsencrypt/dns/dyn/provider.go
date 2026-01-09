package dyn

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/dyn"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

const (
	customerNameFieldID = "dynDnsCustomerName"
	usernameFieldID     = "dynDnsUsername"
	passwordFieldID     = "dynDnsPassword"
)

type Provider struct{}

func (p *Provider) ID() string { return "DYN_DNS" }

func (p *Provider) Name() string { return "Dyn DNS" }

func (p *Provider) DynamicFields() []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          customerNameFieldID,
			Description: "Dyn DNS customer name",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          usernameFieldID,
			Description: "Dyn DNS username",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          passwordFieldID,
			Description: "Dyn DNS password",
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
