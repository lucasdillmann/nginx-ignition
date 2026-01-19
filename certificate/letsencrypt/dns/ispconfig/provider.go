package ispconfig

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/ispconfig"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

const (
	serverURLFieldID          = "ispconfigServerUrl"
	usernameFieldID           = "ispconfigUsername"
	passwordFieldID           = "ispconfigPassword"
	insecureSkipVerifyFieldID = "ispconfigInsecureSkipVerify"
)

type Provider struct{}

func (p *Provider) ID() string {
	return "ISPCONFIG"
}

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateLetsencryptDnsIspconfigName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          serverURLFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsIspconfigServerUrl),
			Required:    true,
			Type:        dynamicfields.URLType,
		},
		{
			ID:          usernameFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsIspconfigUsername),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          passwordFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsIspconfigPassword),
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID: insecureSkipVerifyFieldID,
			Description: i18n.M(
				ctx,
				i18n.K.CertificateLetsencryptDnsIspconfigInsecureSkipVerify,
			),
			Required: false,
			Type:     dynamicfields.BooleanType,
		},
	})
}

func (p *Provider) ChallengeProvider(
	_ context.Context,
	_ []string,
	parameters map[string]any,
) (challenge.Provider, error) {
	serverURL, _ := parameters[serverURLFieldID].(string)
	username, _ := parameters[usernameFieldID].(string)
	password, _ := parameters[passwordFieldID].(string)
	insecureSkipVerify, _ := parameters[insecureSkipVerifyFieldID].(bool)

	cfg := ispconfig.NewDefaultConfig()
	cfg.ServerURL = serverURL
	cfg.Username = username
	cfg.Password = password
	cfg.InsecureSkipVerify = insecureSkipVerify
	cfg.TTL = dns.TTL
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval

	return ispconfig.NewDNSProviderConfig(cfg)
}
