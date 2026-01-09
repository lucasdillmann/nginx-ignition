package lightsail

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/lightsail"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/ptr"
)

const (
	dnsZoneFieldID = "lightsailDnsZone"
	regionFieldID  = "lightsailRegion"
)

type Provider struct{}

func (p *Provider) ID() string { return "LIGHTSAIL" }

func (p *Provider) Name() string { return "AWS Lightsail" }

func (p *Provider) DynamicFields() []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          dnsZoneFieldID,
			Description: "DNS zone name",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          regionFieldID,
			Description: "AWS region",
			HelpText:    ptr.Of("Defaults to us-east-1 when left empty"),
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
