package stackpath

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/stackpath"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

const (
	clientIDFieldID     = "stackPathClientId"
	clientSecretFieldID = "stackPathClientSecret"
	stackIDFieldID      = "stackPathStackId"
)

type Provider struct{}

func (p *Provider) ID() string {
	return "STACKPATH"
}

func (p *Provider) Name() string {
	return "StackPath"
}

func (p *Provider) DynamicFields() []*dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          clientIDFieldID,
			Description: "StackPath client ID",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          clientSecretFieldID,
			Description: "StackPath client secret",
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          stackIDFieldID,
			Description: "StackPath stack ID",
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
	clientID, _ := parameters[clientIDFieldID].(string)
	clientSecret, _ := parameters[clientSecretFieldID].(string)
	stackID, _ := parameters[stackIDFieldID].(string)

	cfg := &stackpath.Config{
		ClientID:           clientID,
		ClientSecret:       clientSecret,
		StackID:            stackID,
		TTL:                dns.TTL,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PollingInterval,
	}

	return stackpath.NewDNSProviderConfig(cfg)
}
