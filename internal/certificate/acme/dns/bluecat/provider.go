package bluecat

import (
	"context"

	"github.com/go-acme/lego/v5/challenge"
	"github.com/go-acme/lego/v5/providers/dns/bluecat"

	"github.com/lucasdillmann/nginx-ignition/internal/certificate/acme/dns"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/dynamicfields"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/i18n"
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
	return i18n.M(ctx, i18n.K.CertificateAcmeDnsBluecatName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          baseURLFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateAcmeDnsBluecatBaseUrl),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          usernameFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateAcmeDnsBluecatUsername),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          passwordFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateAcmeDnsBluecatPassword),
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          configNameFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateAcmeDnsBluecatConfigName),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          dnsViewFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateAcmeDnsBluecatDnsView),
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
