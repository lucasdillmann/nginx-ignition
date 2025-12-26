package httpnet

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/httpnet"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

const (
	apiKeyFieldID   = "httpNetApiKey"
	zoneNameFieldID = "httpNetZoneName"
)

type Provider struct{}

func (p *Provider) ID() string { return "HTTPNET" }

func (p *Provider) Name() string { return "HTTP.net" }

func (p *Provider) DynamicFields() []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          apiKeyFieldID,
			Description: "HTTP.net API key",
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          zoneNameFieldID,
			Description: "HTTP.net zone name",
			Type:        dynamicfields.SingleLineTextType,
		},
	})
}

func (p *Provider) ChallengeProvider(
	_ context.Context,
	_ []string,
	parameters map[string]any,
) (challenge.Provider, error) {
	apiKey, _ := parameters[apiKeyFieldID].(string)
	zoneName, _ := parameters[zoneNameFieldID].(string)

	cfg := &httpnet.Config{
		APIKey:             apiKey,
		ZoneName:           zoneName,
		TTL:                dns.TTL,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PollingInterval,
	}

	return httpnet.NewDNSProviderConfig(cfg)
}
