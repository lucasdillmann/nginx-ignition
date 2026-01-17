package cloudns

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/cloudns"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

const (
	authIDFieldID    = "cloudnsAuthId"
	subAuthIDFieldID = "cloudnsSubAuthId"
	passwordFieldID  = "cloudnsPassword"
)

type Provider struct{}

func (p *Provider) ID() string { return "CLOUDNS" }

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsCloudnsName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          authIDFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsCloudnsAuthId),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          subAuthIDFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsCloudnsSubAuthId),
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          passwordFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsCloudnsPassword),
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
	})
}

func (p *Provider) ChallengeProvider(
	_ context.Context,
	_ []string,
	parameters map[string]any,
) (challenge.Provider, error) {
	authID, _ := parameters[authIDFieldID].(string)
	subAuthID, _ := parameters[subAuthIDFieldID].(string)
	password, _ := parameters[passwordFieldID].(string)

	cfg := cloudns.NewDefaultConfig()
	cfg.AuthID = authID
	cfg.SubAuthID = subAuthID
	cfg.AuthPassword = password
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval
	cfg.TTL = dns.TTL

	return cloudns.NewDNSProviderConfig(cfg)
}
