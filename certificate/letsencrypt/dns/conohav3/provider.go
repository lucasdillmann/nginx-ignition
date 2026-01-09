package conohav3

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/conohav3"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

const (
	userIDFieldID   = "conohaV3UserId"
	passwordFieldID = "conohaV3Password"
	tenantIDFieldID = "conohaV3TenantId"
	regionFieldID   = "conohaV3Region"
)

type Provider struct{}

func (p *Provider) ID() string { return "CONOHA_V3" }

func (p *Provider) Name() string { return "ConoHa v3" }

func (p *Provider) DynamicFields() []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          userIDFieldID,
			Description: "ConoHa user ID",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          passwordFieldID,
			Description: "ConoHa password",
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          tenantIDFieldID,
			Description: "ConoHa tenant ID",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          regionFieldID,
			Description: "ConoHa region",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
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

	cfg := conohav3.NewDefaultConfig()
	cfg.UserID = userName
	cfg.Password = password
	cfg.TenantID = tenantID
	cfg.Region = region
	cfg.TTL = dns.TTL
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval

	return conohav3.NewDNSProviderConfig(cfg)
}
