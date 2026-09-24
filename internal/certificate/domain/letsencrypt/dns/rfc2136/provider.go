package rfc2136

import (
	"context"

	"github.com/go-acme/lego/v5/challenge"
	"github.com/go-acme/lego/v5/providers/dns/dnsupdate"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/dynamicfields"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/i18n"
	dns2 "github.com/lucasdillmann/nginx-ignition/internal/certificate/domain/letsencrypt/dns"
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
	return dns2.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
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

	cfg := dnsupdate.NewDefaultConfig()
	cfg.Nameserver = nameserver
	cfg.TSIGKey = tsigKey
	cfg.TSIGSecret = tsigSecret
	cfg.TSIGAlgorithm = tsigAlgorithm
	cfg.PropagationTimeout = dns2.PropagationTimeout
	cfg.PollingInterval = dns2.PollingInterval
	cfg.TTL = dns2.TTL

	return dnsupdate.NewDNSProviderConfig(cfg)
}
