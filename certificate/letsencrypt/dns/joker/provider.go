package joker

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/joker"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

//nolint:gosec
const (
	dmapi           = "DMAPI"
	svc             = "SVC"
	apiKeyFieldID   = "jokerApiKey"
	usernameFieldID = "jokerUsername"
	passwordFieldID = "jokerPassword"
	apiModeFieldID  = "jokerApiMode"
)

type Provider struct{}

func (p *Provider) ID() string { return "JOKER" }

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateLetsencryptDnsJokerName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:           apiModeFieldID,
			Description:  i18n.M(ctx, i18n.K.CertificateLetsencryptDnsJokerApiMode),
			Type:         dynamicfields.EnumType,
			Required:     true,
			DefaultValue: dmapi,
			EnumOptions: []dynamicfields.EnumOption{
				{
					ID:          dmapi,
					Description: i18n.Static(dmapi),
				},
				{
					ID:          svc,
					Description: i18n.Static(svc),
				},
			},
		},
		{
			ID:          apiKeyFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsJokerApiKey),
			Sensitive:   true,
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
			Conditions: []dynamicfields.Condition{{
				ParentField: apiModeFieldID,
				Value:       dmapi,
			}},
		},
		{
			ID:          usernameFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsJokerUsername),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
			Conditions: []dynamicfields.Condition{{
				ParentField: apiModeFieldID,
				Value:       svc,
			}},
		},
		{
			ID:          passwordFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsJokerPassword),
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
			Conditions: []dynamicfields.Condition{{
				ParentField: apiModeFieldID,
				Value:       svc,
			}},
		},
	})
}

func (p *Provider) ChallengeProvider(
	_ context.Context,
	_ []string,
	parameters map[string]any,
) (challenge.Provider, error) {
	apiKey, _ := parameters[apiKeyFieldID].(string)
	username, _ := parameters[usernameFieldID].(string)
	password, _ := parameters[passwordFieldID].(string)
	apiMode, _ := parameters[apiModeFieldID].(string)

	cfg := joker.NewDefaultConfig()
	cfg.APIKey = apiKey
	cfg.Username = username
	cfg.Password = password
	cfg.APIMode = apiMode
	cfg.TTL = dns.TTL
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval

	return joker.NewDNSProviderConfig(cfg)
}
