package loopia

import (
	"context"

	"github.com/go-acme/lego/v5/challenge"
	"github.com/go-acme/lego/v5/providers/dns/loopia"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/dynamicfields"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/i18n"
	dns2 "github.com/lucasdillmann/nginx-ignition/internal/certificate/domain/letsencrypt/dns"
)

const (
	apiUserFieldID     = "loopiaApiUser"
	apiPasswordFieldID = "loopiaApiPassword"
)

type Provider struct{}

func (p *Provider) ID() string { return "LOOPIA" }

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateLetsencryptDnsLoopiaName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns2.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          apiUserFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsLoopiaApiUser),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          apiPasswordFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsLoopiaApiPassword),
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
	apiUser, _ := parameters[apiUserFieldID].(string)
	apiPassword, _ := parameters[apiPasswordFieldID].(string)

	cfg := loopia.NewDefaultConfig()
	cfg.APIUser = apiUser
	cfg.APIPassword = apiPassword
	cfg.PropagationTimeout = dns2.PropagationTimeout
	cfg.PollingInterval = dns2.PollingInterval
	cfg.TTL = dns2.TTL

	return loopia.NewDNSProviderConfig(cfg)
}
