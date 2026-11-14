package godaddyv3

import (
	"context"

	"github.com/go-acme/lego/v5/challenge"
	"github.com/go-acme/lego/v5/providers/dns/godaddyv3"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/dynamicfields"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/i18n"
	"github.com/lucasdillmann/nginx-ignition/internal/certificate/domain/letsencrypt/dns"
)

//nolint:gosec
const (
	personalAccessTokenFieldID = "goDaddyV3PersonalAccessToken"
	minimumTTL                 = 600
)

type Provider struct{}

func (p *Provider) ID() string {
	return "GODADDY_V3"
}

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateLetsencryptDnsGodaddyv3Name)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          personalAccessTokenFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsGodaddyv3PersonalAccessToken),
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
	})
}

func (p *Provider) ChallengeProvider(
	_ context.Context,
	_ []string,
	parameters map[string]any,
) (challenge.Provider, error) {
	personalAccessToken, _ := parameters[personalAccessTokenFieldID].(string)

	cfg := godaddyv3.NewDefaultConfig()
	cfg.PersonalAccessToken = personalAccessToken
	cfg.TTL = minimumTTL
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval

	return godaddyv3.NewDNSProviderConfig(cfg)
}
