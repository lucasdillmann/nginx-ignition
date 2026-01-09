package inwx

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/inwx"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

//nolint:gosec
const (
	usernameFieldID     = "inwxUsername"
	passwordFieldID     = "inwxPassword"
	sharedSecretFieldID = "inwxSharedSecret"
)

type Provider struct{}

func (p *Provider) ID() string { return "INWX" }

func (p *Provider) Name() string { return "INWX" }

func (p *Provider) DynamicFields() []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          usernameFieldID,
			Description: "INWX username",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          passwordFieldID,
			Description: "INWX password",
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          sharedSecretFieldID,
			Description: "INWX shared secret",
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
	sharedSecret, _ := parameters[sharedSecretFieldID].(string)

	cfg := inwx.NewDefaultConfig()
	cfg.Username = username
	cfg.Password = password
	cfg.SharedSecret = sharedSecret
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval
	cfg.TTL = dns.TTL

	return inwx.NewDNSProviderConfig(cfg)
}
