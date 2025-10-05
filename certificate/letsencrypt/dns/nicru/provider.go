package nicru

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/nicru"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

const (
	usernameFieldID  = "nicRuUsername"
	passwordFieldID  = "nicRuPassword"
	serviceIDFieldID = "nicRuServiceId"
	secretFieldID    = "nicRuSecret"
)

type Provider struct{}

func (p *Provider) ID() string { return "NICRU" }

func (p *Provider) Name() string { return "RU-CENTER" }

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamic_fields.DynamicField{
		{
			ID:          usernameFieldID,
			Description: "RU-CENTER username",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          passwordFieldID,
			Description: "RU-CENTER password",
			Required:    true,
			Sensitive:   true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          serviceIDFieldID,
			Description: "RU-CENTER OAuth2 service ID",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          secretFieldID,
			Description: "RU-CENTER OAuth2 secret",
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
	username, _ := parameters[usernameFieldID].(string)
	password, _ := parameters[passwordFieldID].(string)
	serviceID, _ := parameters[serviceIDFieldID].(string)
	secret, _ := parameters[secretFieldID].(string)

	cfg := &nicru.Config{
		Username:           username,
		Password:           password,
		ServiceID:          serviceID,
		Secret:             secret,
		TTL:                dns.TTL,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PollingInterval,
	}

	return nicru.NewDNSProviderConfig(cfg)
}
