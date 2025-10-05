package joker

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/joker"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

const (
	apiKeyFieldID   = "jokerApiKey"
	usernameFieldID = "jokerUsername"
	passwordFieldID = "jokerPassword"
	apiModeFieldID  = "jokerApiMode"
)

type Provider struct{}

func (p *Provider) ID() string { return "JOKER" }

func (p *Provider) Name() string { return "Joker" }

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamic_fields.DynamicField{
		{
			ID:          apiModeFieldID,
			Description: "Joker API mode",
			Type:        dynamic_fields.EnumType,
			Required:    true,
			EnumOptions: &[]*dynamic_fields.EnumOption{
				{
					ID:          "DMAPI",
					Description: "DMAPI",
				},
				{
					ID:          "SVC",
					Description: "SVC",
				},
			},
		},
		{
			ID:          apiKeyFieldID,
			Description: "Joker API key",
			Sensitive:   true,
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
			Condition: &dynamic_fields.Condition{
				ParentField: apiModeFieldID,
				Value:       "DMAPI",
			},
		},
		{
			ID:          usernameFieldID,
			Description: "Joker username",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
			Condition: &dynamic_fields.Condition{
				ParentField: apiModeFieldID,
				Value:       "SVC",
			},
		},
		{
			ID:          passwordFieldID,
			Description: "Joker password",
			Required:    true,
			Sensitive:   true,
			Type:        dynamic_fields.SingleLineTextType,
			Condition: &dynamic_fields.Condition{
				ParentField: apiModeFieldID,
				Value:       "SVC",
			},
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
		PollingInterval:    dns.PoolingInterval,
	}

	return joker.NewDNSProviderConfig(cfg)
}
