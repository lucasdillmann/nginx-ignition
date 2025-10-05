package otc

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/otc"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

const (
	domainNameFieldID       = "otcDomainName"
	projectNameFieldID      = "otcProjectName"
	userNameFieldID         = "otcUserName"
	passwordFieldID         = "otcPassword"
	identityEndpointFieldID = "otcIdentityEndpoint"
)

type Provider struct{}

func (p *Provider) ID() string { return "OTC" }

func (p *Provider) Name() string { return "Open Telekom Cloud" }

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamic_fields.DynamicField{
		{
			ID:          domainNameFieldID,
			Description: "OTC domain name",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          projectNameFieldID,
			Description: "OTC project name",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          userNameFieldID,
			Description: "OTC username",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          passwordFieldID,
			Description: "OTC password",
			Required:    true,
			Sensitive:   true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          identityEndpointFieldID,
			Description: "OTC identity endpoint",
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
	domainName, _ := parameters[domainNameFieldID].(string)
	projectName, _ := parameters[projectNameFieldID].(string)
	userName, _ := parameters[userNameFieldID].(string)
	password, _ := parameters[passwordFieldID].(string)
	identityEndpoint, _ := parameters[identityEndpointFieldID].(string)

	cfg := &otc.Config{
		DomainName:         domainName,
		ProjectName:        projectName,
		UserName:           userName,
		Password:           password,
		IdentityEndpoint:   identityEndpoint,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PollingInterval,
		SequenceInterval:   dns.SequenceInterval,
	}

	return otc.NewDNSProviderConfig(cfg)
}
