package namecheap

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/namecheap"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

const (
	apiUserFieldID = "namecheapApiUser"
	apiKeyFieldID  = "namecheapApiKey"
)

type Provider struct{}

func (p *Provider) ID() string { return "NAMECHEAP" }

func (p *Provider) Name() string { return "Namecheap" }

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamic_fields.DynamicField{
		{
			ID:          apiUserFieldID,
			Description: "Namecheap API user",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          apiKeyFieldID,
			Description: "Namecheap API key",
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
	apiUser, _ := parameters[apiUserFieldID].(string)
	apiKey, _ := parameters[apiKeyFieldID].(string)

	cfg := &namecheap.Config{
		APIUser:            apiUser,
		APIKey:             apiKey,
		TTL:                dns.TTL,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PollingInterval,
	}

	return namecheap.NewDNSProviderConfig(cfg)
}
