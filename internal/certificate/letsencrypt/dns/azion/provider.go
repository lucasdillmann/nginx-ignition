package azion

import (
	"context"

	"github.com/go-acme/lego/v5/challenge"
	"github.com/go-acme/lego/v5/providers/dns/azion"

	"github.com/lucasdillmann/nginx-ignition/internal/certificate/letsencrypt/dns"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/dynamicfields"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/i18n"
)

const (
	personalTokenFieldID = "azionPersonalToken"
)

type Provider struct{}

func (p *Provider) ID() string { return "AZION" }

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateLetsencryptDnsAzionName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          personalTokenFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsAzionPersonalToken),
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
	personalToken, _ := parameters[personalTokenFieldID].(string)

	cfg := azion.NewDefaultConfig()
	cfg.PersonalToken = personalToken
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval
	cfg.TTL = dns.TTL

	return azion.NewDNSProviderConfig(cfg)
}
