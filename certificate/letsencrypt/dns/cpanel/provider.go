package cpanel

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/cpanel"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

const (
	hostFieldID  = "cpanelHost"
	tokenFieldID = "cpanelToken"
	userFieldID  = "cpanelUsername"
	modeFieldID  = "cpanelMode"
)

type Provider struct{}

func (p *Provider) ID() string { return "CPANEL" }

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsCpanelName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          hostFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsCpanelBaseUrl),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          userFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsCpanelUsername),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          tokenFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsCpanelApiToken),
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          modeFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsCpanelMode),
			HelpText:    i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsCpanelModeHelp),
			Type:        dynamicfields.SingleLineTextType,
		},
	})
}

func (p *Provider) ChallengeProvider(
	_ context.Context,
	_ []string,
	parameters map[string]any,
) (challenge.Provider, error) {
	host, _ := parameters[hostFieldID].(string)
	user, _ := parameters[userFieldID].(string)
	token, _ := parameters[tokenFieldID].(string)
	mode, _ := parameters[modeFieldID].(string)

	if mode == "" {
		mode = "cpanel"
	}

	cfg := cpanel.NewDefaultConfig()
	cfg.BaseURL = host
	cfg.Username = user
	cfg.Token = token
	cfg.Mode = mode
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval

	return cpanel.NewDNSProviderConfig(cfg)
}
