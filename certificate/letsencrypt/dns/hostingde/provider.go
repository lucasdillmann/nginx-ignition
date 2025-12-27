package hostingde

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/hostingde"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

//nolint:gosec
const (
	apiKeyFieldID   = "hostingDeApiKey"
	zoneNameFieldID = "hostingDeZoneName"
)

type Provider struct{}

func (p *Provider) ID() string { return "HOSTINGDE" }

func (p *Provider) Name() string { return "Hosting.de" }

func (p *Provider) DynamicFields() []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          apiKeyFieldID,
			Description: "Hosting.de API key",
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          zoneNameFieldID,
			Description: "Hosting.de zone name",
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

	cfg := &hostingde.Config{
		APIKey:             apiKey,
		ZoneName:           zoneName,
		TTL:                dns.TTL,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PollingInterval,
	}

	return hostingde.NewDNSProviderConfig(cfg)
}
