package rackspace

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/rackspace"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

const (
	userFieldID   = "rackspaceUser"
	apiKeyFieldID = "rackspaceApiKey"
)

type Provider struct{}

func (p *Provider) ID() string { return "RACKSPACE" }

func (p *Provider) Name() string { return "Rackspace" }

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamic_fields.DynamicField{
		{
			ID:          userFieldID,
			Description: "Rackspace username",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          apiKeyFieldID,
			Description: "Rackspace API key",
			Required:    true,
			Sensitive:   true,
			Type:        dynamic_fields.SingleLineTextType,
		},
	})
}

func (p *Provider) ChallengeProvider(
	_ context.Context,
	_ []string,
	parameters map[string]any,
) (challenge.Provider, error) {
	username, _ := parameters[userFieldID].(string)
	apiKey, _ := parameters[apiKeyFieldID].(string)

	cfg := &rackspace.Config{
		APIUser:            username,
		APIKey:             apiKey,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PollingInterval,
		TTL:                dns.TTL,
	}

	return rackspace.NewDNSProviderConfig(cfg)
}
