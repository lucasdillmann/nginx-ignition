package safedns

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/safedns"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

const (
	authTokenFieldID = "safeDnsAuthToken"
)

type Provider struct{}

func (p *Provider) ID() string {
	return "SAFEDNS"
}

func (p *Provider) Name() string {
	return "SafeDNS"
}

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamic_fields.DynamicField{
		{
			ID:          authTokenFieldID,
			Description: "SafeDNS auth token",
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
	authToken, _ := parameters[authTokenFieldID].(string)

	cfg := &safedns.Config{
		AuthToken:          authToken,
		TTL:                dns.TTL,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PollingInterval,
	}

	return safedns.NewDNSProviderConfig(cfg)
}
