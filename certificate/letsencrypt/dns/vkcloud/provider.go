package vkcloud

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/vkcloud"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
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

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamic_fields.DynamicField{
		{
			ID:          usernameFieldID,
			Description: "VK Cloud username",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          passwordFieldID,
			Description: "VK Cloud password",
			Required:    true,
			Sensitive:   true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          projectIDFieldID,
			Description: "VK Cloud project ID",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          dnsEndpointFieldID,
			Description: "VK Cloud DNS endpoint",
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          identityEndpointFieldID,
			Description: "VK Cloud identity endpoint",
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          domainNameFieldID,
			Description: "VK Cloud domain name",
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
		PollingInterval:    dns.PoolingInterval,
	}

	return vkcloud.NewDNSProviderConfig(cfg)
}
