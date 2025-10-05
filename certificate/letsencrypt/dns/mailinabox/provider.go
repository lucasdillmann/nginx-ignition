package mailinabox

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/mailinabox"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

const (
	emailFieldID    = "mailInABoxEmail"
	passwordFieldID = "mailInABoxPassword"
)

type Provider struct{}

func (p *Provider) ID() string { return "MAIL_IN_A_BOX" }

func (p *Provider) Name() string { return "Mail-in-a-Box" }

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamic_fields.DynamicField{
		{
			ID:          emailFieldID,
			Description: "Mail-in-a-Box email",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          passwordFieldID,
			Description: "Mail-in-a-Box password",
			Required:    true,
			Sensitive:   true,
			Type:        dynamic_fields.SingleLineTextType,
		},
	})
}

func (p *Provider) ChallengeProvider(_ context.Context, _ []string, parameters map[string]any) (challenge.Provider, error) {
	email, _ := parameters[emailFieldID].(string)
	password, _ := parameters[passwordFieldID].(string)

	cfg := &mailinabox.Config{
		Email:              email,
		Password:           password,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PoolingInterval,
	}

	return mailinabox.NewDNSProviderConfig(cfg)
}
