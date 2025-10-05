package selectelv2

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/selectelv2"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

const (
	usernameFieldID  = "selectelv2Username"
	passwordFieldID  = "selectelv2Password"
	projectIDFieldID = "selectelv2ProjectId"
	accountFieldID   = "selectelv2Account"
	regionFieldID    = "selectelv2Region"
)

type Provider struct{}

func (p *Provider) ID() string {
	return "SELECTELV2"
}

func (p *Provider) Name() string {
	return "Selectel v2"
}

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamic_fields.DynamicField{
		{
			ID:          usernameFieldID,
			Description: "Selectel v2 username",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          passwordFieldID,
			Description: "Selectel v2 password",
			Required:    true,
			Sensitive:   true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          projectIDFieldID,
			Description: "Selectel v2 project ID",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          accountFieldID,
			Description: "Selectel v2 account name",
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          regionFieldID,
			Description: "Selectel v2 region",
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
	projectID, _ := parameters[projectIDFieldID].(string)
	account, _ := parameters[accountFieldID].(string)
	region, _ := parameters[regionFieldID].(string)

	cfg := &selectelv2.Config{
		Username:           username,
		Password:           password,
		ProjectID:          projectID,
		DomainName:         account,
		AuthRegion:         region,
		TTL:                dns.TTL,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PoolingInterval,
	}

	return selectelv2.NewDNSProviderConfig(cfg)
}
