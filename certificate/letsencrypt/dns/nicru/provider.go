package nicru

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/nicru"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
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

func (p *Provider) Name() string { return "RU-CENTER" }

func (p *Provider) DynamicFields() []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          usernameFieldID,
			Description: "RU-CENTER username",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          passwordFieldID,
			Description: "RU-CENTER password",
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          serviceIDFieldID,
			Description: "RU-CENTER OAuth2 service ID",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          secretFieldID,
			Description: "RU-CENTER OAuth2 secret",
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
