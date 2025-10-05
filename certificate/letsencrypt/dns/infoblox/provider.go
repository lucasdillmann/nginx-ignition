package infoblox

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/infoblox"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

const (
	hostFieldID          = "infobloxHost"
	portFieldID          = "infobloxPort"
	usernameFieldID      = "infobloxUsername"
	passwordFieldID      = "infobloxPassword"
	dnsViewFieldID       = "infobloxDnsView"
	wapiVersionFieldID   = "infobloxWapiVersion"
	sslVerifyFieldID     = "infobloxSslVerify"
	caCertificateFieldID = "infobloxCaCertificate"
)

type Provider struct{}

func (p *Provider) ID() string { return "INFOBLOX" }

func (p *Provider) Name() string { return "Infoblox" }

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamic_fields.DynamicField{
		{
			ID:          hostFieldID,
			Description: "Infoblox grid manager host",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          portFieldID,
			Description: "Infoblox grid manager port",
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          usernameFieldID,
			Description: "Infoblox username",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          passwordFieldID,
			Description: "Infoblox password",
			Required:    true,
			Sensitive:   true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          dnsViewFieldID,
			Description: "Infoblox DNS view",
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          wapiVersionFieldID,
			Description: "Infoblox WAPI version",
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          sslVerifyFieldID,
			Description: "Verify SSL certificate",
			Type:        dynamic_fields.BooleanType,
		},
		{
			ID:          caCertificateFieldID,
			Description: "CA certificate path (PEM encoded)",
			Type:        dynamic_fields.SingleLineTextType,
		},
	})
}

func (p *Provider) ChallengeProvider(
	_ context.Context,
	_ []string,
	parameters map[string]any,
) (challenge.Provider, error) {
	host, _ := parameters[hostFieldID].(string)
	port, _ := parameters[portFieldID].(string)
	username, _ := parameters[usernameFieldID].(string)
	password, _ := parameters[passwordFieldID].(string)
	dnsView, _ := parameters[dnsViewFieldID].(string)
	wapiVersion, _ := parameters[wapiVersionFieldID].(string)
	sslVerify, _ := parameters[sslVerifyFieldID].(bool)
	caCertificate, _ := parameters[caCertificateFieldID].(string)

	cfg := &infoblox.Config{
		Host:               host,
		Port:               port,
		Username:           username,
		Password:           password,
		DNSView:            dnsView,
		WapiVersion:        wapiVersion,
		SSLVerify:          sslVerify,
		CACertificate:      caCertificate,
		TTL:                dns.TTL,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PollingInterval,
	}

	return infoblox.NewDNSProviderConfig(cfg)
}
