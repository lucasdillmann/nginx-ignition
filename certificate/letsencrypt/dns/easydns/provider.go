package easydns

import (
	"context"
	"net/url"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/easydns"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

//nolint:gosec
const (
	endpointFieldID = "easyDnsEndpoint"
	tokenFieldID    = "easyDnsToken"
	keyFieldID      = "easyDnsKey"
)

type Provider struct{}

func (p *Provider) ID() string { return "EASYDNS" }

func (p *Provider) Name() string { return "EasyDNS" }

func (p *Provider) DynamicFields() []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          tokenFieldID,
			Description: "EasyDNS API token",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          keyFieldID,
			Description: "EasyDNS API key",
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          endpointFieldID,
			Description: "EasyDNS API endpoint",
			Type:        dynamicfields.SingleLineTextType,
		},
	})
}

func (p *Provider) ChallengeProvider(
	_ context.Context,
	_ []string,
	parameters map[string]any,
) (challenge.Provider, error) {
	token, _ := parameters[tokenFieldID].(string)
	apiKey, _ := parameters[keyFieldID].(string)
	endpoint, _ := parameters[endpointFieldID].(string)

	cfg := easydns.NewDefaultConfig()
	cfg.Token = token
	cfg.Key = apiKey
	cfg.TTL = dns.TTL
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.SequenceInterval = dns.SequenceInterval
	cfg.PollingInterval = dns.PollingInterval

	if endpoint != "" {
		if parsedValue, err := url.Parse(endpoint); err == nil {
			cfg.Endpoint = parsedValue
		}
	}

	return easydns.NewDNSProviderConfig(cfg)
}
