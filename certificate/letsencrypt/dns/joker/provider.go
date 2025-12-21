package joker

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/joker"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/ptr"
)

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

func (p *Provider) Name() string { return "Joker" }

func (p *Provider) DynamicFields() []*dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:           apiModeFieldID,
			Description:  "Joker API mode",
			Type:         dynamicfields.EnumType,
			Required:     true,
			DefaultValue: ptr.Of(dmapi),
			EnumOptions: &[]*dynamicfields.EnumOption{
				{
					ID:          dmapi,
					Description: dmapi,
				},
				{
					ID:          svc,
					Description: svc,
				},
			},
		},
		{
			ID:          apiKeyFieldID,
			Description: "Joker API key",
			Sensitive:   true,
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
			Conditions: &[]dynamicfields.Condition{{
				ParentField: apiModeFieldID,
				Value:       dmapi,
			}},
		},
		{
			ID:          usernameFieldID,
			Description: "Joker username",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
			Conditions: &[]dynamicfields.Condition{{
				ParentField: apiModeFieldID,
				Value:       svc,
			}},
		},
		{
			ID:          passwordFieldID,
			Description: "Joker password",
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
			Conditions: &[]dynamicfields.Condition{{
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

	cfg := &joker.Config{
		APIKey:             apiKey,
		Username:           username,
		Password:           password,
		APIMode:            apiMode,
		TTL:                dns.TTL,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PollingInterval,
	}

	return joker.NewDNSProviderConfig(cfg)
}
