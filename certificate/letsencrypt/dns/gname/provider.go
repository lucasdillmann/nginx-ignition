package gname

import (
	"context"

	"github.com/go-acme/lego/v5/challenge"
	"github.com/go-acme/lego/v5/providers/dns/gname"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

const (
	appIDFieldID  = "gnameAppId"
	appKeyFieldID = "gnameAppKey"
)

type Provider struct{}

func (p *Provider) ID() string {
	return "GNAME"
}

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateLetsencryptDnsGnameName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          appIDFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsGnameAppId),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          appKeyFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsGnameAppKey),
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
	appID, _ := parameters[appIDFieldID].(string)
	appKey, _ := parameters[appKeyFieldID].(string)

	cfg := gname.NewDefaultConfig()
	cfg.AppID = appID
	cfg.AppKey = appKey
	cfg.TTL = dns.TTL
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval

	return gname.NewDNSProviderConfig(cfg)
}
