package rfc2136

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/rfc2136"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

//nolint:gosec
const (
	nameserverFieldID    = "rfc2136Nameserver"
	tsigKeyFieldID       = "rfc2136TsigKey"
	tsigSecretFieldID    = "rfc2136TsigSecret"
	tsigAlgorithmFieldID = "rfc2136TsigAlgorithm"
)

type Provider struct{}

func (p *Provider) ID() string { return "RFC2136" }

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateLetsencryptDnsRfc2136Name)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID: nameserverFieldID,
			Description: i18n.M(
				ctx,
				i18n.K.CertificateLetsencryptDnsRfc2136NameserverAddress,
			),
			HelpText: i18n.M(
				ctx,
				i18n.K.CertificateLetsencryptDnsRfc2136NameserverAddressHelp,
			),
			Required: true,
			Type:     dynamicfields.SingleLineTextType,
		},
		{
			ID:          tsigKeyFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsRfc2136TsigKeyName),
			HelpText:    i18n.M(ctx, i18n.K.CertificateLetsencryptDnsRfc2136TsigKeyNameHelp),
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          tsigSecretFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsRfc2136TsigSecretKey),
			HelpText: i18n.M(
				ctx,
				i18n.K.CertificateLetsencryptDnsRfc2136TsigSecretKeyHelp,
			),
			Sensitive: true,
			Type:      dynamicfields.SingleLineTextType,
		},
		{
			ID:          tsigAlgorithmFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsRfc2136TsigAlgorithm),
			HelpText: i18n.M(
				ctx,
				i18n.K.CertificateLetsencryptDnsRfc2136TsigAlgorithmHelp,
			),
			Type: dynamicfields.SingleLineTextType,
		},
	})
}

func (p *Provider) ChallengeProvider(
	_ context.Context,
	_ []string,
	parameters map[string]any,
) (challenge.Provider, error) {
	nameserver, _ := parameters[nameserverFieldID].(string)
	tsigKey, _ := parameters[tsigKeyFieldID].(string)
	tsigSecret, _ := parameters[tsigSecretFieldID].(string)
	tsigAlgorithm, _ := parameters[tsigAlgorithmFieldID].(string)

	cfg := rfc2136.NewDefaultConfig()
	cfg.Nameserver = nameserver
	cfg.TSIGKey = tsigKey
	cfg.TSIGSecret = tsigSecret
	cfg.TSIGAlgorithm = tsigAlgorithm
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval
	cfg.TTL = dns.TTL

	return rfc2136.NewDNSProviderConfig(cfg)
}
