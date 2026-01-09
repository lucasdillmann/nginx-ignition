package rackspace

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/rackspace"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

//nolint:gosec
const (
	userFieldID   = "rackspaceUser"
	apiKeyFieldID = "rackspaceApiKey"
)

type Provider struct{}

func (p *Provider) ID() string { return "RACKSPACE" }

func (p *Provider) Name() string { return "Rackspace" }

func (p *Provider) DynamicFields() []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          userFieldID,
			Description: "Rackspace username",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          apiKeyFieldID,
			Description: "Rackspace API key",
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
	username, _ := parameters[userFieldID].(string)
	apiKey, _ := parameters[apiKeyFieldID].(string)

	cfg := rackspace.NewDefaultConfig()
	cfg.APIUser = username
	cfg.APIKey = apiKey
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval
	cfg.TTL = dns.TTL

	return rackspace.NewDNSProviderConfig(cfg)
}
