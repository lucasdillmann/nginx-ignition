package efficientip

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/efficientip"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
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

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateLetsencryptDnsEfficientipName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          usernameFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsEfficientipUsername),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          passwordFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsEfficientipPassword),
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          hostnameFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsEfficientipHostname),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID: dnsNameFieldID,
			Description: i18n.M(
				ctx,
				i18n.K.CertificateLetsencryptDnsEfficientipDnsServerName,
			),
			Required: true,
			Type:     dynamicfields.SingleLineTextType,
		},
		{
			ID:          viewNameFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsEfficientipDnsViewName),
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID: insecureSkipVerifyFieldID,
			Description: i18n.M(
				ctx,
				i18n.K.CertificateLetsencryptDnsEfficientipSkipTlsVerify,
			),
			Type: dynamicfields.BooleanType,
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
