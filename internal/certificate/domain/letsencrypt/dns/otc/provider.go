package otc

import (
	"context"

	"github.com/go-acme/lego/v5/challenge"
	"github.com/go-acme/lego/v5/providers/dns/otc"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/dynamicfields"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/i18n"
	dns2 "github.com/lucasdillmann/nginx-ignition/internal/certificate/domain/letsencrypt/dns"
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

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateLetsencryptDnsOtcName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns2.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          domainNameFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsOtcDomainName),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          projectNameFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsOtcProjectName),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          userNameFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsOtcUsername),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          passwordFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsOtcPassword),
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          identityEndpointFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsOtcIdentityEndpoint),
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

	cfg := otc.NewDefaultConfig()
	cfg.DomainName = domainName
	cfg.ProjectName = projectName
	cfg.UserName = userName
	cfg.Password = password
	cfg.IdentityEndpoint = identityEndpoint
	cfg.PropagationTimeout = dns2.PropagationTimeout
	cfg.PollingInterval = dns2.PollingInterval
	cfg.SequenceInterval = dns2.SequenceInterval

	return otc.NewDNSProviderConfig(cfg)
}
