package azion

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/azion"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

const (
	personalTokenFieldID = "azionPersonalToken"
)

type Provider struct{}

func (p *Provider) ID() string { return "AZION" }

func (p *Provider) Name() string { return "Azion" }

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamic_fields.DynamicField{
		{
			ID:          personalTokenFieldID,
			Description: "Azion Personal Token",
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
	personalToken, _ := parameters[personalTokenFieldID].(string)

	cfg := &azion.Config{
		PersonalToken:      personalToken,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PoolingInterval,
		TTL:                dns.TTL,
	}

	return azion.NewDNSProviderConfig(cfg)
}
