package lightsail

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/lightsail"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

const (
	dnsZoneFieldID = "lightsailDnsZone"
	regionFieldID  = "lightsailRegion"
)

type Provider struct{}

func (p *Provider) ID() string { return "LIGHTSAIL" }

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsLightsailName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          dnsZoneFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsLightsailDnsZoneName),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          regionFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsLightsailAwsRegion),
			HelpText:    i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsLightsailAwsRegionHelp),
			Type:        dynamicfields.SingleLineTextType,
		},
	})
}

func (p *Provider) ChallengeProvider(
	_ context.Context,
	_ []string,
	parameters map[string]any,
) (challenge.Provider, error) {
	dnsZone, _ := parameters[dnsZoneFieldID].(string)
	region, _ := parameters[regionFieldID].(string)

	if region == "" {
		region = "us-east-1"
	}

	cfg := lightsail.NewDefaultConfig()
	cfg.DNSZone = dnsZone
	cfg.Region = region
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval

	return lightsail.NewDNSProviderConfig(cfg)
}
