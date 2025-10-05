package iwantmyname

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/iwantmyname"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

const (
	usernameFieldID = "iWantMyNameUsername"
	passwordFieldID = "iWantMyNamePassword"
)

type Provider struct{}

func (p *Provider) ID() string { return "IWANTMYNAME" }

func (p *Provider) Name() string { return "iwantmyname" }

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamic_fields.DynamicField{
		{
			ID:          usernameFieldID,
			Description: "iwantmyname username",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          passwordFieldID,
			Description: "iwantmyname password",
			Required:    true,
			Sensitive:   true,
			Type:        dynamic_fields.SingleLineTextType,
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

	cfg := &iwantmyname.Config{
		Username:           username,
		Password:           password,
		TTL:                dns.TTL,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PoolingInterval,
	}

	return iwantmyname.NewDNSProviderConfig(cfg)
}
