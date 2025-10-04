package hostingde

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/hostingde"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

const (
	apiKeyFieldID   = "hostingDeApiKey"
	zoneNameFieldID = "hostingDeZoneName"
)

type Provider struct{}

func (p *Provider) ID() string { return "HOSTINGDE" }

func (p *Provider) Name() string { return "Hosting.de" }

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamic_fields.DynamicField{
		{
			ID:          apiKeyFieldID,
			Description: "Hosting.de API key",
			Required:    true,
			Sensitive:   true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          zoneNameFieldID,
			Description: "Hosting.de zone name",
			Required:    false,
			Type:        dynamic_fields.SingleLineTextType,
		},
	})
}

func (p *Provider) ChallengeProvider(_ context.Context, _ []string, parameters map[string]any) (challenge.Provider, error) {
	apiKey, _ := parameters[apiKeyFieldID].(string)
	zoneName, _ := parameters[zoneNameFieldID].(string)

	cfg := &hostingde.Config{
		APIKey:             apiKey,
		ZoneName:           zoneName,
		TTL:                dns.TTL,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PoolingInterval,
	}

	return hostingde.NewDNSProviderConfig(cfg)
}
