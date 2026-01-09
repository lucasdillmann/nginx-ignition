package namedotcom

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/namedotcom"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

const (
	usernameFieldID = "nameDotComUsername"
	apiTokenFieldID = "nameDotComApiToken"
	serverFieldID   = "nameDotComServer"
)

type Provider struct{}

func (p *Provider) ID() string { return "NAMEDOTCOM" }

func (p *Provider) Name() string { return "Name.com" }

func (p *Provider) DynamicFields() []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          usernameFieldID,
			Description: "Name.com username",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          apiTokenFieldID,
			Description: "Name.com API token",
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          serverFieldID,
			Description: "Name.com server",
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
	apiToken, _ := parameters[apiTokenFieldID].(string)
	server, _ := parameters[serverFieldID].(string)

	cfg := namedotcom.NewDefaultConfig()
	cfg.Username = username
	cfg.APIToken = apiToken
	cfg.Server = server
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval
	cfg.TTL = dns.TTL

	return namedotcom.NewDNSProviderConfig(cfg)
}
