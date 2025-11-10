package f5xc

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/f5xc"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

const (
	apiTokenFieldID   = "f5xcApiToken"
	tenantNameFieldID = "f5xcTenantName"
	groupNameFieldID  = "f5xcGroupName"
)

type Provider struct{}

func (p *Provider) ID() string { return "F5_XC" }

func (p *Provider) Name() string { return "F5 XC" }

func (p *Provider) DynamicFields() []*dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          apiTokenFieldID,
			Description: "F5 XC API token",
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          tenantNameFieldID,
			Description: "F5 XC tenant name",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          groupNameFieldID,
			Description: "F5 XC group name",
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
	apiToken, _ := parameters[apiTokenFieldID].(string)
	tenantName, _ := parameters[tenantNameFieldID].(string)
	groupName, _ := parameters[groupNameFieldID].(string)

	cfg := &f5xc.Config{
		APIToken:           apiToken,
		TenantName:         tenantName,
		GroupName:          groupName,
		TTL:                dns.TTL,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PollingInterval,
	}

	return f5xc.NewDNSProviderConfig(cfg)
}
