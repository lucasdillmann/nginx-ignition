package efficientip

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/efficientip"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

const (
	usernameFieldID           = "efficientIpUsername"
	passwordFieldID           = "efficientIpPassword"
	hostnameFieldID           = "efficientIpHostname"
	dnsNameFieldID            = "efficientIpDnsName"
	viewNameFieldID           = "efficientIpViewName"
	insecureSkipVerifyFieldID = "efficientIpInsecureSkipVerify"
)

type Provider struct{}

func (p *Provider) ID() string { return "EFFICIENTIP" }

func (p *Provider) Name() string { return "EfficientIP" }

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamic_fields.DynamicField{
		{
			ID:          usernameFieldID,
			Description: "EfficientIP username",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          passwordFieldID,
			Description: "EfficientIP password",
			Required:    true,
			Sensitive:   true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          hostnameFieldID,
			Description: "EfficientIP hostname",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          dnsNameFieldID,
			Description: "DNS server name",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          viewNameFieldID,
			Description: "DNS view name",
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          insecureSkipVerifyFieldID,
			Description: "Skip TLS certificate verification",
			Type:        dynamic_fields.BooleanType,
		},
	})
}

func (p *Provider) ChallengeProvider(
	_ context.Context,
	_ []string,
	parameters map[string]any,
) (challenge.Provider, error) {
	username, _ := parameters[usernameFieldID].(string)
	password, _ := parameters[passwordFieldID].(string)
	hostname, _ := parameters[hostnameFieldID].(string)
	dnsName, _ := parameters[dnsNameFieldID].(string)
	viewName, _ := parameters[viewNameFieldID].(string)
	insecureSkipVerify, _ := parameters[insecureSkipVerifyFieldID].(bool)

	cfg := &efficientip.Config{
		Username:           username,
		Password:           password,
		Hostname:           hostname,
		DNSName:            dnsName,
		ViewName:           viewName,
		InsecureSkipVerify: insecureSkipVerify,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PollingInterval,
	}

	return efficientip.NewDNSProviderConfig(cfg)
}
