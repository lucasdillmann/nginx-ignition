package lightsail

import (
	"context"

	"github.com/go-acme/lego/v5/challenge"
	"github.com/go-acme/lego/v5/providers/dns/lightsail"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/dynamicfields"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/i18n"
	dns2 "github.com/lucasdillmann/nginx-ignition/internal/certificate/domain/letsencrypt/dns"
)

const (
	dnsZoneFieldID = "lightsailDnsZone"
	regionFieldID  = "lightsailRegion"
)

type Provider struct{}

func (p *Provider) ID() string { return "LIGHTSAIL" }

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateLetsencryptDnsLightsailName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns2.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          dnsZoneFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsLightsailDnsZoneName),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          regionFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsLightsailAwsRegion),
			HelpText:    i18n.M(ctx, i18n.K.CertificateLetsencryptDnsLightsailAwsRegionHelp),
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
	cfg.PropagationTimeout = dns2.PropagationTimeout
	cfg.PollingInterval = dns2.PollingInterval

	return lightsail.NewDNSProviderConfig(cfg)
}
