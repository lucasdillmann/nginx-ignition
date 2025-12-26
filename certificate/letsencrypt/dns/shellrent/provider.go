package shellrent

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/shellrent"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

const (
	usernameFieldID = "shellrentUsername"
	tokenFieldID    = "shellrentToken"
)

type Provider struct{}

func (p *Provider) ID() string {
	return "SHELLRENT"
}

func (p *Provider) Name() string {
	return "Shellrent"
}

func (p *Provider) DynamicFields() []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          usernameFieldID,
			Description: "Shellrent username",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          tokenFieldID,
			Description: "Shellrent token",
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
	token, _ := parameters[tokenFieldID].(string)

	cfg := &shellrent.Config{
		Username:           username,
		Token:              token,
		TTL:                dns.TTL,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PollingInterval,
	}

	return shellrent.NewDNSProviderConfig(cfg)
}
