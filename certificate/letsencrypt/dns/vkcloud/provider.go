package vkcloud

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/vkcloud"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

const (
	usernameFieldID         = "vkCloudUsername"
	passwordFieldID         = "vkCloudPassword"
	projectIDFieldID        = "vkCloudProjectId"
	dnsEndpointFieldID      = "vkCloudDnsEndpoint"
	identityEndpointFieldID = "vkCloudIdentityEndpoint"
	domainNameFieldID       = "vkCloudDomainName"
)

type Provider struct{}

func (p *Provider) ID() string {
	return "VK_CLOUD"
}

func (p *Provider) Name() string {
	return "VK Cloud"
}

func (p *Provider) DynamicFields() []*dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          usernameFieldID,
			Description: "VK Cloud username",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          passwordFieldID,
			Description: "VK Cloud password",
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          projectIDFieldID,
			Description: "VK Cloud project ID",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          dnsEndpointFieldID,
			Description: "VK Cloud DNS endpoint",
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          identityEndpointFieldID,
			Description: "VK Cloud identity endpoint",
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          domainNameFieldID,
			Description: "VK Cloud domain name",
			Type:        dynamicfields.SingleLineTextType,
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
	dnsEndpoint, _ := parameters[dnsEndpointFieldID].(string)
	identityEndpoint, _ := parameters[identityEndpointFieldID].(string)
	domainName, _ := parameters[domainNameFieldID].(string)

	cfg := &vkcloud.Config{
		Username:           username,
		Password:           password,
		ProjectID:          projectID,
		DNSEndpoint:        dnsEndpoint,
		IdentityEndpoint:   identityEndpoint,
		DomainName:         domainName,
		TTL:                dns.TTL,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PollingInterval,
	}

	return vkcloud.NewDNSProviderConfig(cfg)
}
