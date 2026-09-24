package cloudns

import (
	"context"

	"github.com/go-acme/lego/v5/challenge"
	"github.com/go-acme/lego/v5/providers/dns/cloudns"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/dynamicfields"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/i18n"
	dns2 "github.com/lucasdillmann/nginx-ignition/internal/certificate/domain/letsencrypt/dns"
)

const (
	authIDFieldID    = "cloudnsAuthId"
	subAuthIDFieldID = "cloudnsSubAuthId"
	passwordFieldID  = "cloudnsPassword"
)

type Provider struct{}

func (p *Provider) ID() string { return "CLOUDNS" }

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateLetsencryptDnsCloudnsName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns2.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          authIDFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsCloudnsAuthId),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          subAuthIDFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsCloudnsSubAuthId),
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          passwordFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsCloudnsPassword),
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
	cfg.PropagationTimeout = dns2.PropagationTimeout
	cfg.PollingInterval = dns2.PollingInterval
	cfg.TTL = dns2.TTL

	return cloudns.NewDNSProviderConfig(cfg)
}
