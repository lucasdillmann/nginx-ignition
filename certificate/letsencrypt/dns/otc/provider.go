package otc

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/otc"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
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

func (p *Provider) DynamicFields() []*dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          domainNameFieldID,
			Description: "OTC domain name",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          projectNameFieldID,
			Description: "OTC project name",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          userNameFieldID,
			Description: "OTC username",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          passwordFieldID,
			Description: "OTC password",
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          identityEndpointFieldID,
			Description: "OTC identity endpoint",
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
