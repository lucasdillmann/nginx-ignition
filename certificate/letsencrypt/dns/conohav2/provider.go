package conohav2

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/conoha"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

const (
	regionFieldID   = "conoHaV2Region"
	tenantIDFieldID = "conoHaV2TenantID"
	usernameFieldID = "conoHaV2Username"
	passwordFieldID = "conoHaV2Password"
)

type Provider struct{}

func (p *Provider) ID() string { return "CONOHA_V2" }

func (p *Provider) Name() string { return "ConoHa v2" }

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
		PollingInterval:    dns.PollingInterval,
		TTL:                dns.TTL,
	}

	return conoha.NewDNSProviderConfig(cfg)
}
