package gandiv5

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/gandiv5"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

const (
	personalAccessTokenFieldID = "gandiPersonalAccessToken"
)

type Provider struct{}

func (p *Provider) ID() string { return "GANDI_V5" }

func (p *Provider) Name() string { return "Gandi v5 (LiveDNS)" }

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamic_fields.DynamicField{
		{
			ID:          personalAccessTokenFieldID,
			Description: "Gandi Personal Access Token (PAT)",
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
	token, _ := parameters[personalAccessTokenFieldID].(string)

	cfg := &gandiv5.Config{
		PersonalAccessToken: token,
		PropagationTimeout:  dns.PropagationTimeout,
		PollingInterval:     dns.PoolingInterval,
		TTL:                 dns.TTL,
	}

	return gandiv5.NewDNSProviderConfig(cfg)
}
