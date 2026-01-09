package safedns

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/safedns"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

//nolint:gosec
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

func (p *Provider) DynamicFields() []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          authTokenFieldID,
			Description: "SafeDNS auth token",
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
	authToken, _ := parameters[authTokenFieldID].(string)

	cfg := safedns.NewDefaultConfig()
	cfg.AuthToken = authToken
	cfg.TTL = dns.TTL
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval

	return safedns.NewDNSProviderConfig(cfg)
}
