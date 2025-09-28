package ovh

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/ovh"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

const (
	endpointFieldID       = "ovhEndpoint"
	applicationKeyFieldID = "ovhApplicationKey"
	applicationSecFieldID = "ovhApplicationSecret"
	consumerKeyFieldID    = "ovhConsumerKey"
)

type Provider struct{}

func (p *Provider) ID() string { return "OVH" }

func (p *Provider) Name() string { return "OVH" }

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamic_fields.DynamicField{
		{
			ID:          endpointFieldID,
			Description: "OVH API endpoint",
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          applicationKeyFieldID,
			Description: "OVH application key",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          applicationSecFieldID,
			Description: "OVH application secret",
			Required:    true,
			Sensitive:   true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          consumerKeyFieldID,
			Description: "OVH consumer key",
			Required:    true,
			Sensitive:   true,
			Type:        dynamic_fields.SingleLineTextType,
		},
	})
}

func (p *Provider) ChallengeProvider(_ context.Context, _ []string, parameters map[string]any) (challenge.Provider, error) {
	endpoint, _ := parameters[endpointFieldID].(string)
	appKey, _ := parameters[applicationKeyFieldID].(string)
	appSecret, _ := parameters[applicationSecFieldID].(string)
	consumerKey, _ := parameters[consumerKeyFieldID].(string)

	cfg := &ovh.Config{
		APIEndpoint:        endpoint,
		ApplicationKey:     appKey,
		ApplicationSecret:  appSecret,
		ConsumerKey:        consumerKey,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PoolingInterval,
		TTL:                dns.TTL,
	}

	return ovh.NewDNSProviderConfig(cfg)
}
