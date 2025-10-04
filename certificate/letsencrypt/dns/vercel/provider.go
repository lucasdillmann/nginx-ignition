package vercel

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/vercel"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

const (
	tokenFieldID = "vercelToken"
	teamFieldID  = "vercelTeamId"
)

type Provider struct{}

func (p *Provider) ID() string { return "VERCEL" }

func (p *Provider) Name() string { return "Vercel" }

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamic_fields.DynamicField{
		{
			ID:          tokenFieldID,
			Description: "Vercel token",
			Required:    true,
			Sensitive:   true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          teamFieldID,
			Description: "Vercel team ID",
			Required:    false,
			Type:        dynamic_fields.SingleLineTextType,
		},
	})
}

func (p *Provider) ChallengeProvider(_ context.Context, _ []string, parameters map[string]any) (challenge.Provider, error) {
	teamId, _ := parameters[teamFieldID].(string)
	token, _ := parameters[tokenFieldID].(string)

	cfg := &vercel.Config{
		AuthToken:          token,
		TeamID:             teamId,
		TTL:                dns.TTL,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PoolingInterval,
	}

	return vercel.NewDNSProviderConfig(cfg)
}
