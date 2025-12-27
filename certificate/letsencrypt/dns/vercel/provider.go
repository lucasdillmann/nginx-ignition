package vercel

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/vercel"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

//nolint:gosec
const (
	tokenFieldID = "vercelToken"
	teamFieldID  = "vercelTeamId"
)

type Provider struct{}

func (p *Provider) ID() string { return "VERCEL" }

func (p *Provider) Name() string { return "Vercel" }

func (p *Provider) DynamicFields() []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          tokenFieldID,
			Description: "Vercel token",
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          teamFieldID,
			Description: "Vercel team ID",
			Type:        dynamicfields.SingleLineTextType,
		},
	})
}

func (p *Provider) ChallengeProvider(
	_ context.Context,
	_ []string,
	parameters map[string]any,
) (challenge.Provider, error) {
	teamId, _ := parameters[teamFieldID].(string)
	token, _ := parameters[tokenFieldID].(string)

	cfg := &vercel.Config{
		AuthToken:          token,
		TeamID:             teamId,
		TTL:                dns.TTL,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PollingInterval,
	}

	return vercel.NewDNSProviderConfig(cfg)
}
