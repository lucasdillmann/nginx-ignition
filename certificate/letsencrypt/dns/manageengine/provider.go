package manageengine

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/manageengine"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

const (
	clientIDFieldID     = "manageEngineClientId"
	clientSecretFieldID = "manageEngineClientSecret"
)

type Provider struct{}

func (p *Provider) ID() string { return "MANAGEENGINE" }

func (p *Provider) Name() string { return "ManageEngine CloudDNS" }

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamic_fields.DynamicField{
		{
			ID:          clientIDFieldID,
			Description: "ManageEngine CloudDNS Client ID",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          clientSecretFieldID,
			Description: "ManageEngine CloudDNS Client secret",
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
	clientID, _ := parameters[clientIDFieldID].(string)
	clientSecret, _ := parameters[clientSecretFieldID].(string)

	cfg := &manageengine.Config{
		ClientID:           clientID,
		ClientSecret:       clientSecret,
		TTL:                dns.TTL,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PollingInterval,
	}

	return manageengine.NewDNSProviderConfig(cfg)
}
