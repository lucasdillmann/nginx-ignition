package loopia

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/loopia"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

const (
	apiUserFieldID     = "loopiaApiUser"
	apiPasswordFieldID = "loopiaApiPassword"
)

type Provider struct{}

func (p *Provider) ID() string { return "LOOPIA" }

func (p *Provider) Name() string { return "Loopia" }

func (p *Provider) DynamicFields() []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          apiUserFieldID,
			Description: "Loopia API user",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          apiPasswordFieldID,
			Description: "Loopia API password",
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
	apiUser, _ := parameters[apiUserFieldID].(string)
	apiPassword, _ := parameters[apiPasswordFieldID].(string)

	cfg := loopia.NewDefaultConfig()
	cfg.APIUser = apiUser
	cfg.APIPassword = apiPassword
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval
	cfg.TTL = dns.TTL

	return loopia.NewDNSProviderConfig(cfg)
}
