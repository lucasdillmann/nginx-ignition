package liquidweb

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/liquidweb"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

const (
	usernameFieldID = "liquidWebUsername"
	passwordFieldID = "liquidWebPassword"
	zoneFieldID     = "liquidWebZone"
)

type Provider struct{}

func (p *Provider) ID() string { return "LIQUIDWEB" }

func (p *Provider) Name() string { return "LiquidWeb" }

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamic_fields.DynamicField{
		{
			ID:          usernameFieldID,
			Description: "LiquidWeb username",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          passwordFieldID,
			Description: "LiquidWeb password",
			Required:    true,
			Sensitive:   true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          zoneFieldID,
			Description: "LiquidWeb zone",
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
	zone, _ := parameters[zoneFieldID].(string)

	cfg := &liquidweb.Config{
		Username:           username,
		Password:           password,
		Zone:               zone,
		TTL:                dns.TTL,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PoolingInterval,
	}

	return liquidweb.NewDNSProviderConfig(cfg)
}
