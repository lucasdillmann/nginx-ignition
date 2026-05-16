package dnsla

import (
	"context"

	"github.com/go-acme/lego/v5/challenge"
	"github.com/go-acme/lego/v5/providers/dns/dnsla"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

const (
	apiIDFieldID     = "dnslaApiId"
	apiSecretFieldID = "dnslaApiSecret"
)

type Provider struct{}

func (p *Provider) ID() string {
	return "DNSLA"
}

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateLetsencryptDnsDnslaName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          apiIDFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsDnslaApiId),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          apiSecretFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsDnslaApiSecret),
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
	apiID, _ := parameters[apiIDFieldID].(string)
	apiSecret, _ := parameters[apiSecretFieldID].(string)

	cfg := dnsla.NewDefaultConfig()
	cfg.APIID = apiID
	cfg.APISecret = apiSecret
	cfg.TTL = dns.TTL
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval

	return dnsla.NewDNSProviderConfig(cfg)
}
