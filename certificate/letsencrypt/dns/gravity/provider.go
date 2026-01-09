package gravity

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/gravity"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

const (
	usernameFieldID  = "gravityUsername"
	passwordFieldID  = "gravityPassword"
	serverURLFieldID = "gravityServerUrl"
)

type Provider struct{}

func (p *Provider) ID() string {
	return "GRAVITY"
}

func (p *Provider) Name() string {
	return "Gravity"
}

func (p *Provider) DynamicFields() []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          usernameFieldID,
			Description: "Gravity username",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          passwordFieldID,
			Description: "Gravity password",
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          serverURLFieldID,
			Description: "Gravity server URL",
			Required:    true,
			Type:        dynamicfields.URLType,
		},
	})
}

func (p *Provider) ChallengeProvider(
	_ context.Context,
	_ []string,
	parameters map[string]any,
) (challenge.Provider, error) {
	username, _ := parameters[usernameFieldID].(string)
	password, _ := parameters[passwordFieldID].(string)
	serverURL, _ := parameters[serverURLFieldID].(string)

	cfg := gravity.NewDefaultConfig()
	cfg.Username = username
	cfg.Password = password
	cfg.ServerURL = serverURL
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval

	return gravity.NewDNSProviderConfig(cfg)
}
