package conohav3

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/conohav3"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

const (
	userIDFieldID   = "conohaUserId"
	passwordFieldID = "conohaPassword"
	tenantIDFieldID = "conohaTenantId"
	regionFieldID   = "conohaRegion"
)

type Provider struct{}

func (p *Provider) ID() string { return "CONOHA" }

func (p *Provider) Name() string { return "ConoHa" }

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamic_fields.DynamicField{
		{
			ID:          userIDFieldID,
			Description: "ConoHa user ID",
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
		{
			ID:          tenantIDFieldID,
			Description: "ConoHa tenant ID",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          regionFieldID,
			Description: "ConoHa region",
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
	userName, _ := parameters[userIDFieldID].(string)
	password, _ := parameters[passwordFieldID].(string)
	tenantID, _ := parameters[tenantIDFieldID].(string)
	region, _ := parameters[regionFieldID].(string)

	cfg := &conohav3.Config{
		UserID:             userName,
		Password:           password,
		TenantID:           tenantID,
		Region:             region,
		TTL:                dns.TTL,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PoolingInterval,
	}

	return conohav3.NewDNSProviderConfig(cfg)
}
