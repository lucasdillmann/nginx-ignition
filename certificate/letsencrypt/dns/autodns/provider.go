package autodns

import (
	"context"
	"fmt"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/autodns"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

const (
	usernameFieldID = "autoDnsUsername"
	passwordFieldID = "autoDnsPassword"
	contextFieldID  = "autoDnsContext"
)

type Provider struct{}

func (p *Provider) ID() string { return "AUTODNS" }

func (p *Provider) Name() string { return "AutoDNS" }

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamic_fields.DynamicField{
		{
			ID:          usernameFieldID,
			Description: "AutoDNS username",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          passwordFieldID,
			Description: "AutoDNS password",
			Required:    true,
			Sensitive:   true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          contextFieldID,
			Description: "AutoDNS context",
			Required:    true,
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
	contextStr, _ := parameters[contextFieldID].(string)

	var contextInt int
	if contextStr != "" {
		_, err := fmt.Sscanf(contextStr, "%d", &contextInt)
		if err != nil {
			return nil, err
		}
	}

	cfg := &autodns.Config{
		Username:           username,
		Password:           password,
		Context:            contextInt,
		TTL:                dns.TTL,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PoolingInterval,
	}

	return autodns.NewDNSProviderConfig(cfg)
}
