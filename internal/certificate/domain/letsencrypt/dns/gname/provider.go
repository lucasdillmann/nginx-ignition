package gname

import (
	"context"

	"github.com/go-acme/lego/v5/challenge"
	"github.com/go-acme/lego/v5/providers/dns/gname"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/dynamicfields"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/i18n"
	dns2 "github.com/lucasdillmann/nginx-ignition/internal/certificate/domain/letsencrypt/dns"
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
	return dns2.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
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
	cfg.TTL = dns2.TTL
	cfg.PropagationTimeout = dns2.PropagationTimeout
	cfg.PollingInterval = dns2.PollingInterval

	return gname.NewDNSProviderConfig(cfg)
}
