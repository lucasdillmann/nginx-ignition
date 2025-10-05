package allinkl

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/allinkl"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

const (
	loginFieldID    = "allInklLogin"
	passwordFieldID = "allInklPassword"
)

type Provider struct{}

func (p *Provider) ID() string { return "ALL_INKL" }

func (p *Provider) Name() string { return "ALL-INKL" }

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamic_fields.DynamicField{
		{
			ID:          loginFieldID,
			Description: "ALL-INKL login",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          passwordFieldID,
			Description: "ALL-INKL password",
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
	login, _ := parameters[loginFieldID].(string)
	password, _ := parameters[passwordFieldID].(string)

	cfg := &allinkl.Config{
		Login:              login,
		Password:           password,
		TTL:                dns.TTL,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PoolingInterval,
	}

	return allinkl.NewDNSProviderConfig(cfg)
}
