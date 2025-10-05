package dyndnsfree

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/dyndnsfree"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

const (
	usernameFieldID = "dynDnsFreeUsername"
	passwordFieldID = "dynDnsFreePassword"
)

type Provider struct{}

func (p *Provider) ID() string { return "DYN_DNS_FREE" }

func (p *Provider) Name() string { return "Dyn DNS Free" }

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamic_fields.DynamicField{
		{
			ID:          usernameFieldID,
			Description: "DynDNSFree username",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          passwordFieldID,
			Description: "DynDNSFree password",
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
	user, _ := parameters[usernameFieldID].(string)
	pass, _ := parameters[passwordFieldID].(string)

	cfg := &dyndnsfree.Config{
		Username:           user,
		Password:           pass,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PoolingInterval,
	}

	return dyndnsfree.NewDNSProviderConfig(cfg)
}
