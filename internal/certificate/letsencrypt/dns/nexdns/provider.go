package nexdns

import (
	"context"

	"github.com/go-acme/lego/v5/challenge"
	"github.com/go-acme/lego/v5/providers/dns/nexdns"

	"nginx-ignition/internal/certificate/letsencrypt/dns"
	"nginx-ignition/internal/core/common/dynamicfields"
	"nginx-ignition/internal/core/common/i18n"
)

//nolint:gosec
const (
	apiTokenFieldID = "nexdnsApiToken"
)

type Provider struct{}

func (p *Provider) ID() string {
	return "NEXDNS"
}

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateLetsencryptDnsNexdnsName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          apiTokenFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsNexdnsApiToken),
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
	apiToken, _ := parameters[apiTokenFieldID].(string)

	cfg := nexdns.NewDefaultConfig()
	cfg.APIToken = apiToken
	cfg.TTL = dns.TTL
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval

	return nexdns.NewDNSProviderConfig(cfg)
}
