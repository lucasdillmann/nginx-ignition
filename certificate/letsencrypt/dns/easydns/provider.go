package easydns

import (
	"context"
	"net/url"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/easydns"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

const (
	endpointFieldID = "easyDnsEndpoint"
	tokenFieldID    = "easyDnsToken"
	keyFieldID      = "easyDnsKey"
)

type Provider struct{}

func (p *Provider) ID() string { return "EASYDNS" }

func (p *Provider) Name() string { return "EasyDNS" }

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamic_fields.DynamicField{
		{
			ID:          tokenFieldID,
			Description: "EasyDNS API token",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          keyFieldID,
			Description: "EasyDNS API key",
			Required:    true,
			Sensitive:   true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          endpointFieldID,
			Description: "EasyDNS API endpoint",
			Type:        dynamic_fields.SingleLineTextType,
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

	cfg := &easydns.Config{
		Token:              token,
		Key:                apiKey,
		TTL:                dns.TTL,
		PropagationTimeout: dns.PropagationTimeout,
		SequenceInterval:   dns.SequenceInterval,
		PollingInterval:    dns.PollingInterval,
	}

	if endpoint != "" {
		if parsedValue, err := url.Parse(endpoint); err == nil {
			cfg.Endpoint = parsedValue
		}
	}

	return easydns.NewDNSProviderConfig(cfg)
}
