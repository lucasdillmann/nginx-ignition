package versio

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/versio"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

const (
	usernameFieldID = "versioUsername"
	passwordFieldID = "versioPassword"
)

type Provider struct{}

func (p *Provider) ID() string {
	return "VERSIO"
}

func (p *Provider) Name() string {
	return "Versio"
}

func (p *Provider) DynamicFields() []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          usernameFieldID,
			Description: "Versio username",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          passwordFieldID,
			Description: "Versio password",
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

	cfg := &versio.Config{
		Username:           username,
		Password:           password,
		TTL:                dns.TTL,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PollingInterval,
		SequenceInterval:   dns.SequenceInterval,
	}

	return versio.NewDNSProviderConfig(cfg)
}
