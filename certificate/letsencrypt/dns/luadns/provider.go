package luadns

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/luadns"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

const (
	apiUsernameFieldID = "luaDnsApiUsername"
	apiTokenFieldID    = "luaDnsApiToken"
)

type Provider struct{}

func (p *Provider) ID() string { return "LUADNS" }

func (p *Provider) Name() string { return "LuaDNS" }

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamic_fields.DynamicField{
		{
			ID:          apiUsernameFieldID,
			Description: "LuaDNS API username",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          apiTokenFieldID,
			Description: "LuaDNS API token",
			Required:    true,
			Sensitive:   true,
			Type:        dynamic_fields.SingleLineTextType,
		},
	})
}

func (p *Provider) ChallengeProvider(_ context.Context, _ []string, parameters map[string]any) (challenge.Provider, error) {
	apiUsername, _ := parameters[apiUsernameFieldID].(string)
	apiToken, _ := parameters[apiTokenFieldID].(string)

	cfg := &luadns.Config{
		APIUsername:        apiUsername,
		APIToken:           apiToken,
		TTL:                dns.TTL,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PoolingInterval,
	}

	return luadns.NewDNSProviderConfig(cfg)
}
