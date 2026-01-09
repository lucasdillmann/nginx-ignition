package ispconfig

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/ispconfig"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
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

func (p *Provider) Name() string {
	return "ISPConfig 3"
}

func (p *Provider) DynamicFields() []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          serverURLFieldID,
			Description: "ISPConfig 3 server URL",
			Required:    true,
			Type:        dynamicfields.URLType,
		},
		{
			ID:          usernameFieldID,
			Description: "ISPConfig 3 username",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          passwordFieldID,
			Description: "ISPConfig 3 password",
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          insecureSkipVerifyFieldID,
			Description: "Insecure skip verify (TLS)",
			Required:    false,
			Type:        dynamicfields.BooleanType,
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
