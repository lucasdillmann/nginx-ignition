package vercel

import (
	"context"

	"github.com/go-acme/lego/v5/challenge"
	"github.com/go-acme/lego/v5/providers/dns/vercel"

	"github.com/lucasdillmann/nginx-ignition/internal/certificate/letsencrypt/dns"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/dynamicfields"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/i18n"
)

//nolint:gosec
const (
	tokenFieldID = "vercelToken"
	teamFieldID  = "vercelTeamId"
)

type Provider struct{}

func (p *Provider) ID() string { return "VERCEL" }

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateLetsencryptDnsVercelName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          tokenFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsVercelToken),
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          teamFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsVercelTeamId),
			Type:        dynamicfields.SingleLineTextType,
		},
	})
}

func (p *Provider) ChallengeProvider(
	_ context.Context,
	_ []string,
	parameters map[string]any,
) (challenge.Provider, error) {
	teamID, _ := parameters[teamFieldID].(string)
	token, _ := parameters[tokenFieldID].(string)

	cfg := vercel.NewDefaultConfig()
	cfg.AuthToken = token
	cfg.TeamID = teamID
	cfg.TTL = dns.TTL
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval

	return vercel.NewDNSProviderConfig(cfg)
}
