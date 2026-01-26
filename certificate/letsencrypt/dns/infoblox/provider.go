package infoblox

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/infoblox"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
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

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateLetsencryptDnsInfobloxName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          hostFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsInfobloxGridManagerHost),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          portFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsInfobloxGridManagerPort),
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          usernameFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsInfobloxUsername),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          passwordFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsInfobloxPassword),
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          dnsViewFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsInfobloxDnsView),
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          wapiVersionFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsInfobloxWapiVersion),
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          sslVerifyFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsInfobloxVerifySsl),
			Type:        dynamicfields.BooleanType,
		},
		{
			ID: caCertificateFieldID,
			Description: i18n.M(
				ctx,
				i18n.K.CertificateLetsencryptDnsInfobloxCaCertificatePath,
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
	host, _ := parameters[hostFieldID].(string)
	port, _ := parameters[portFieldID].(string)
	username, _ := parameters[usernameFieldID].(string)
	password, _ := parameters[passwordFieldID].(string)
	dnsView, _ := parameters[dnsViewFieldID].(string)
	wapiVersion, _ := parameters[wapiVersionFieldID].(string)
	sslVerify, _ := parameters[sslVerifyFieldID].(bool)
	caCertificate, _ := parameters[caCertificateFieldID].(string)

	cfg := infoblox.NewDefaultConfig()
	cfg.Host = host
	cfg.Port = port
	cfg.Username = username
	cfg.Password = password
	cfg.DNSView = dnsView
	cfg.WapiVersion = wapiVersion
	cfg.SSLVerify = sslVerify
	cfg.CACertificate = caCertificate
	cfg.TTL = dns.TTL
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval

	return infoblox.NewDNSProviderConfig(cfg)
}
