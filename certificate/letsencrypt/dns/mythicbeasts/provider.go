package mythicbeasts

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/mythicbeasts"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

const (
	usernameFieldID = "mythicBeastsUsername"
	passwordFieldID = "mythicBeastsPassword"
)

type Provider struct{}

func (p *Provider) ID() string { return "MYTHIC_BEASTS" }

func (p *Provider) Name() string { return "Mythic Beasts" }

func (p *Provider) DynamicFields() []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          usernameFieldID,
			Description: "Mythic Beasts username",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          passwordFieldID,
			Description: "Mythic Beasts password",
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

	cfg := &mythicbeasts.Config{
		UserName:           username,
		Password:           password,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PollingInterval,
		TTL:                dns.TTL,
	}

	return mythicbeasts.NewDNSProviderConfig(cfg)
}
