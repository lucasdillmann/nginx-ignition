package efficientip

import (
	"context"

	"github.com/go-acme/lego/v5/challenge"
	"github.com/go-acme/lego/v5/providers/dns/efficientip"

	"github.com/lucasdillmann/nginx-ignition/internal/certificate/acme/dns"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/dynamicfields"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/i18n"
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
	return i18n.M(ctx, i18n.K.CertificateAcmeDnsEfficientipName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          usernameFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateAcmeDnsEfficientipUsername),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          passwordFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateAcmeDnsEfficientipPassword),
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          hostnameFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateAcmeDnsEfficientipHostname),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID: dnsNameFieldID,
			Description: i18n.M(
				ctx,
				i18n.K.CertificateAcmeDnsEfficientipDnsServerName,
			),
			Required: true,
			Type:     dynamicfields.SingleLineTextType,
		},
		{
			ID:          viewNameFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateAcmeDnsEfficientipDnsViewName),
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID: insecureSkipVerifyFieldID,
			Description: i18n.M(
				ctx,
				i18n.K.CertificateAcmeDnsEfficientipSkipTlsVerify,
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
