package inwx

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/inwx"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

const (
	usernameFieldID     = "inwxUsername"
	passwordFieldID     = "inwxPassword"
	sharedSecretFieldID = "inwxSharedSecret"
)

type Provider struct{}

func (p *Provider) ID() string { return "INWX" }

func (p *Provider) Name() string { return "INWX" }

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamic_fields.DynamicField{
		{
			ID:          usernameFieldID,
			Description: "INWX username",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          passwordFieldID,
			Description: "INWX password",
			Required:    true,
			Sensitive:   true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          sharedSecretFieldID,
			Description: "INWX shared secret",
			Sensitive:   true,
			Type:        dynamic_fields.SingleLineTextType,
		},
	})
}

func (p *Provider) ChallengeProvider(_ context.Context, _ []string, parameters map[string]any) (challenge.Provider, error) {
	username, _ := parameters[usernameFieldID].(string)
	password, _ := parameters[passwordFieldID].(string)
	sharedSecret, _ := parameters[sharedSecretFieldID].(string)

	cfg := &inwx.Config{
		Username:           username,
		Password:           password,
		SharedSecret:       sharedSecret,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PoolingInterval,
		TTL:                dns.TTL,
	}

	return inwx.NewDNSProviderConfig(cfg)
}
