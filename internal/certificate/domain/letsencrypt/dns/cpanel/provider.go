package cpanel

import (
	"context"

	"github.com/go-acme/lego/v5/challenge"
	"github.com/go-acme/lego/v5/providers/dns/cpanel"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/dynamicfields"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/i18n"
	dns2 "github.com/lucasdillmann/nginx-ignition/internal/certificate/domain/letsencrypt/dns"
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
	return i18n.M(ctx, i18n.K.CertificateLetsencryptDnsCpanelName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns2.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          hostFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsCpanelBaseUrl),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          userFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsCpanelUsername),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          tokenFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsCpanelApiToken),
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          modeFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsCpanelMode),
			HelpText:    i18n.M(ctx, i18n.K.CertificateLetsencryptDnsCpanelModeHelp),
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
	cfg.PropagationTimeout = dns2.PropagationTimeout
	cfg.PollingInterval = dns2.PollingInterval

	return cpanel.NewDNSProviderConfig(cfg)
}
