package selfhostde

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/selfhostde"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

//nolint:gosec
const (
	usernameFieldID = "selfHostDeUsername"
	passwordFieldID = "selfHostDePassword"
)

type Provider struct{}

func (p *Provider) ID() string {
	return "SELFHOSTDE"
}

func (p *Provider) Name() string {
	return "SelfHost.de"
}

func (p *Provider) DynamicFields() []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          usernameFieldID,
			Description: "SelfHost.de username",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          passwordFieldID,
			Description: "SelfHost.de password",
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

	cfg := &selfhostde.Config{
		Username:           username,
		Password:           password,
		TTL:                dns.TTL,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PollingInterval,
	}

	return selfhostde.NewDNSProviderConfig(cfg)
}
