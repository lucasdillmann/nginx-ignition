package namedotcom

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/namedotcom"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

const (
	usernameFieldID = "nameDotComUsername"
	apiTokenFieldID = "nameDotComApiToken"
	serverFieldID   = "nameDotComServer"
)

type Provider struct{}

func (p *Provider) ID() string { return "NAMEDOTCOM" }

func (p *Provider) Name() string { return "Name.com" }

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamic_fields.DynamicField{
		{
			ID:          usernameFieldID,
			Description: "Name.com username",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          apiTokenFieldID,
			Description: "Name.com API token",
			Required:    true,
			Sensitive:   true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          serverFieldID,
			Description: "Name.com server",
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
	apiToken, _ := parameters[apiTokenFieldID].(string)
	server, _ := parameters[serverFieldID].(string)

	cfg := &namedotcom.Config{
		Username:           username,
		APIToken:           apiToken,
		Server:             server,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PoolingInterval,
		TTL:                dns.TTL,
	}

	return namedotcom.NewDNSProviderConfig(cfg)
}
