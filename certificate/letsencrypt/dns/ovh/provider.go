package ovh

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/ovh"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

//nolint:gosec
const (
	endpointFieldID       = "ovhEndpoint"
	applicationKeyFieldID = "ovhApplicationKey"
	applicationSecFieldID = "ovhApplicationSecret"
	consumerKeyFieldID    = "ovhConsumerKey"
)

type Provider struct{}

func (p *Provider) ID() string { return "OVH" }

func (p *Provider) Name() string { return "OVH" }

func (p *Provider) DynamicFields() []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          endpointFieldID,
			Description: "OVH API endpoint",
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          applicationKeyFieldID,
			Description: "OVH application key",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          applicationSecFieldID,
			Description: "OVH application secret",
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          consumerKeyFieldID,
			Description: "OVH consumer key",
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
	endpoint, _ := parameters[endpointFieldID].(string)
	appKey, _ := parameters[applicationKeyFieldID].(string)
	appSecret, _ := parameters[applicationSecFieldID].(string)
	consumerKey, _ := parameters[consumerKeyFieldID].(string)

	cfg := ovh.NewDefaultConfig()
	cfg.APIEndpoint = endpoint
	cfg.ApplicationKey = appKey
	cfg.ApplicationSecret = appSecret
	cfg.ConsumerKey = consumerKey
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval
	cfg.TTL = dns.TTL

	return ovh.NewDNSProviderConfig(cfg)
}
