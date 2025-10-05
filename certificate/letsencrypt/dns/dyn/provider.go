package dyn

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/dyn"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

const (
	customerNameFieldID = "dynDnsCustomerName"
	usernameFieldID     = "dynDnsUsername"
	passwordFieldID     = "dynDnsPassword"
)

type Provider struct{}

func (p *Provider) ID() string { return "DYN_DNS" }

func (p *Provider) Name() string { return "Dyn DNS" }

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamic_fields.DynamicField{
		{
			ID:          customerNameFieldID,
			Description: "Dyn DNS customer name",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          usernameFieldID,
			Description: "Dyn DNS username",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          passwordFieldID,
			Description: "Dyn DNS password",
			Required:    true,
			Sensitive:   true,
			Type:        dynamic_fields.SingleLineTextType,
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

	cfg := &dyn.Config{
		CustomerName:       customerName,
		UserName:           username,
		Password:           password,
		TTL:                dns.TTL,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PollingInterval,
	}

	return dyn.NewDNSProviderConfig(cfg)
}
