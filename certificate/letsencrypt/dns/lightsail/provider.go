package lightsail

import (
	"context"

	"github.com/aws/smithy-go/ptr"
	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/lightsail"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

const (
	dnsZoneFieldID = "lightsailDnsZone"
	regionFieldID  = "lightsailRegion"
)

type Provider struct{}

func (p *Provider) ID() string { return "LIGHTSAIL" }

func (p *Provider) Name() string { return "AWS Lightsail" }

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamic_fields.DynamicField{
		{
			ID:          dnsZoneFieldID,
			Description: "DNS zone name",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          regionFieldID,
			Description: "AWS region",
			HelpText:    ptr.String("Defaults to us-east-1 when left empty"),
			Type:        dynamic_fields.SingleLineTextType,
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

	cfg := &lightsail.Config{
		DNSZone:            dnsZone,
		Region:             region,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PollingInterval,
	}

	return lightsail.NewDNSProviderConfig(cfg)
}
