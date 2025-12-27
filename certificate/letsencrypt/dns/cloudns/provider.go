package cloudns

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/cloudns"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

const (
	authIDFieldID    = "cloudnsAuthId"
	subAuthIDFieldID = "cloudnsSubAuthId"
	passwordFieldID  = "cloudnsPassword"
)

type Provider struct{}

func (p *Provider) ID() string { return "CLOUDNS" }

func (p *Provider) Name() string { return "ClouDNS" }

func (p *Provider) DynamicFields() []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          authIDFieldID,
			Description: "ClouDNS auth ID",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          subAuthIDFieldID,
			Description: "ClouDNS sub auth ID",
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          passwordFieldID,
			Description: "ClouDNS password",
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
	authID, _ := parameters[authIDFieldID].(string)
	subAuthID, _ := parameters[subAuthIDFieldID].(string)
	password, _ := parameters[passwordFieldID].(string)

	cfg := &cloudns.Config{
		AuthID:             authID,
		SubAuthID:          subAuthID,
		AuthPassword:       password,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PollingInterval,
		TTL:                dns.TTL,
	}

	return cloudns.NewDNSProviderConfig(cfg)
}
