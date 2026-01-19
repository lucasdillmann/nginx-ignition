package bluecat

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/bluecat"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

const (
	baseURLFieldID    = "blueCatBaseURL"
	usernameFieldID   = "blueCatUsername"
	passwordFieldID   = "blueCatPassword"
	configNameFieldID = "blueCatConfigName"
	dnsViewFieldID    = "blueCatDNSView"
)

type Provider struct{}

func (p *Provider) ID() string { return "BLUECAT" }

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateLetsencryptDnsBluecatName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          baseURLFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsBluecatBaseUrl),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          usernameFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsBluecatUsername),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          passwordFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsBluecatPassword),
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          configNameFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsBluecatConfigName),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          dnsViewFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsBluecatDnsView),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
	})
}

func (p *Provider) ChallengeProvider(
	_ context.Context,
	_ []string,
	parameters map[string]any,
) (challenge.Provider, error) {
	baseURL, _ := parameters[baseURLFieldID].(string)
	username, _ := parameters[usernameFieldID].(string)
	password, _ := parameters[passwordFieldID].(string)
	configName, _ := parameters[configNameFieldID].(string)
	dnsView, _ := parameters[dnsViewFieldID].(string)

	cfg := bluecat.NewDefaultConfig()
	cfg.BaseURL = baseURL
	cfg.UserName = username
	cfg.Password = password
	cfg.ConfigName = configName
	cfg.DNSView = dnsView
	cfg.TTL = dns.TTL
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval

	return bluecat.NewDNSProviderConfig(cfg)
}
