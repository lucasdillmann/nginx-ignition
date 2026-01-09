package luadns

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/luadns"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

//nolint:gosec
const (
	apiUsernameFieldID = "luaDnsApiUsername"
	apiTokenFieldID    = "luaDnsApiToken"
)

type Provider struct{}

func (p *Provider) ID() string { return "LUADNS" }

func (p *Provider) Name() string { return "LuaDNS" }

func (p *Provider) DynamicFields() []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          apiUsernameFieldID,
			Description: "LuaDNS API username",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          apiTokenFieldID,
			Description: "LuaDNS API token",
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
	apiUsername, _ := parameters[apiUsernameFieldID].(string)
	apiToken, _ := parameters[apiTokenFieldID].(string)

	cfg := luadns.NewDefaultConfig()
	cfg.APIUsername = apiUsername
	cfg.APIToken = apiToken
	cfg.TTL = dns.TTL
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval

	return luadns.NewDNSProviderConfig(cfg)
}
