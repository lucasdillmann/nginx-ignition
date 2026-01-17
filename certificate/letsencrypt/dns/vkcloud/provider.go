package vkcloud

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/vkcloud"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
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

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsVkcloudName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          usernameFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsVkcloudUsername),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          passwordFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsVkcloudPassword),
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          projectIDFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsVkcloudProjectId),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          dnsEndpointFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsVkcloudDnsEndpoint),
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          identityEndpointFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsVkcloudIdentityEndpoint),
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          domainNameFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsVkcloudDomainName),
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

	cfg := vkcloud.NewDefaultConfig()
	cfg.Username = username
	cfg.Password = password
	cfg.ProjectID = projectID
	cfg.DNSEndpoint = dnsEndpoint
	cfg.IdentityEndpoint = identityEndpoint
	cfg.DomainName = domainName
	cfg.TTL = dns.TTL
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval

	return vkcloud.NewDNSProviderConfig(cfg)
}
