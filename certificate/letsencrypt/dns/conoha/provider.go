package conoha

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/conoha"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

const (
	regionFieldID   = "conoHaRegion"
	tenantIDFieldID = "conoHaTenantID"
	usernameFieldID = "conoHaUsername"
	passwordFieldID = "conoHaPassword"
)

type Provider struct{}

func (p *Provider) ID() string { return "CONOHA" }

func (p *Provider) Name() string { return "ConoHa" }

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamic_fields.DynamicField{
		{
			ID:          regionFieldID,
			Description: "ConoHa region",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          tenantIDFieldID,
			Description: "ConoHa tenant ID",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          usernameFieldID,
			Description: "ConoHa username",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          passwordFieldID,
			Description: "ConoHa password",
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
	region, _ := parameters[regionFieldID].(string)
	tenantID, _ := parameters[tenantIDFieldID].(string)
	username, _ := parameters[usernameFieldID].(string)
	password, _ := parameters[passwordFieldID].(string)

	cfg := &conoha.Config{
		Region:             region,
		TenantID:           tenantID,
		Username:           username,
		Password:           password,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PoolingInterval,
		TTL:                dns.TTL,
	}

	return conoha.NewDNSProviderConfig(cfg)
}
