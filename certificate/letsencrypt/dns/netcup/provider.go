package netcup

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/netcup"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

const (
	customerFieldID    = "netcupCustomer"
	apiKeyFieldID      = "netcupApiKey"
	apiPasswordFieldID = "netcupApiPassword"
)

type Provider struct{}

func (p *Provider) ID() string { return "NETCUP" }

func (p *Provider) Name() string { return "Netcup" }

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamic_fields.DynamicField{
		{
			ID:          customerFieldID,
			Description: "Netcup customer",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          apiKeyFieldID,
			Description: "Netcup API key",
			Required:    true,
			Sensitive:   true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          apiPasswordFieldID,
			Description: "Netcup API password",
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
	customerNumber, _ := parameters[customerFieldID].(string)
	apiKey, _ := parameters[apiKeyFieldID].(string)
	apiPassword, _ := parameters[apiPasswordFieldID].(string)

	cfg := &netcup.Config{
		Customer:           customerNumber,
		Key:                apiKey,
		Password:           apiPassword,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PoolingInterval,
		TTL:                dns.TTL,
	}

	return netcup.NewDNSProviderConfig(cfg)
}
