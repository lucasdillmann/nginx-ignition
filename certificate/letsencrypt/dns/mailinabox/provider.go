package mailinabox

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/mailinabox"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

const (
	baseURLFieldID  = "mailInABoxBaseUrl"
	emailFieldID    = "mailInABoxEmail"
	passwordFieldID = "mailInABoxPassword"
)

type Provider struct{}

func (p *Provider) ID() string { return "MAIL_IN_A_BOX" }

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateLetsencryptDnsMailinaboxName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          baseURLFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsMailinaboxBaseUrl),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          emailFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsMailinaboxEmail),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          passwordFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsMailinaboxPassword),
			Required:    true,
			Sensitive:   true,
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
	email, _ := parameters[emailFieldID].(string)
	password, _ := parameters[passwordFieldID].(string)

	cfg := mailinabox.NewDefaultConfig()
	cfg.BaseURL = baseURL
	cfg.Email = email
	cfg.Password = password
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval

	return mailinabox.NewDNSProviderConfig(cfg)
}
