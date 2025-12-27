package hyperone

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/hyperone"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

const (
	apiEndpointFieldID      = "hyperOneApiEndpoint"
	locationIDFieldID       = "hyperOneLocationId"
	passportLocationFieldID = "hyperOnePassportLocation"
)

type Provider struct{}

func (p *Provider) ID() string { return "HYPERONE" }

func (p *Provider) Name() string { return "HyperOne" }

func (p *Provider) DynamicFields() []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          apiEndpointFieldID,
			Description: "HyperOne API endpoint",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          locationIDFieldID,
			Description: "HyperOne location ID",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          passportLocationFieldID,
			Description: "HyperOne passport file location",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
	})
}

func (p *Provider) ChallengeProvider(
	_ context.Context,
	_ []string,
	parameters map[string]any,
) (challenge.Provider, error) {
	apiEndpoint, _ := parameters[apiEndpointFieldID].(string)
	locationID, _ := parameters[locationIDFieldID].(string)
	passportLocation, _ := parameters[passportLocationFieldID].(string)

	cfg := &hyperone.Config{
		APIEndpoint:        apiEndpoint,
		LocationID:         locationID,
		PassportLocation:   passportLocation,
		TTL:                dns.TTL,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PollingInterval,
	}

	return hyperone.NewDNSProviderConfig(cfg)
}
