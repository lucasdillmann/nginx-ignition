package efficientip

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/efficientip"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
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

func (p *Provider) DynamicFields() []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          usernameFieldID,
			Description: "EfficientIP username",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          passwordFieldID,
			Description: "EfficientIP password",
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          hostnameFieldID,
			Description: "EfficientIP hostname",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          dnsNameFieldID,
			Description: "DNS server name",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          viewNameFieldID,
			Description: "DNS view name",
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          insecureSkipVerifyFieldID,
			Description: "Skip TLS certificate verification",
			Type:        dynamicfields.BooleanType,
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

	cfg := efficientip.NewDefaultConfig()
	cfg.Username = username
	cfg.Password = password
	cfg.Hostname = hostname
	cfg.DNSName = dnsName
	cfg.ViewName = viewName
	cfg.InsecureSkipVerify = insecureSkipVerify
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval

	return efficientip.NewDNSProviderConfig(cfg)
}
