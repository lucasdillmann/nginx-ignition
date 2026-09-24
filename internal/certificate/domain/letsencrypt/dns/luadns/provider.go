package luadns

import (
	"context"

	"github.com/go-acme/lego/v5/challenge"
	"github.com/go-acme/lego/v5/providers/dns/luadns"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/dynamicfields"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/i18n"
	dns2 "github.com/lucasdillmann/nginx-ignition/internal/certificate/domain/letsencrypt/dns"
)

//nolint:gosec
const (
	apiUsernameFieldID = "luaDnsApiUsername"
	apiTokenFieldID    = "luaDnsApiToken"
)

type Provider struct{}

func (p *Provider) ID() string { return "LUADNS" }

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateLetsencryptDnsLuadnsName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns2.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          apiUsernameFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsLuadnsApiUsername),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          apiTokenFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsLuadnsApiToken),
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
	apiUsername, _ := parameters[apiUsernameFieldID].(string)
	apiToken, _ := parameters[apiTokenFieldID].(string)

	cfg := luadns.NewDefaultConfig()
	cfg.APIUsername = apiUsername
	cfg.APIToken = apiToken
	cfg.TTL = dns2.TTL
	cfg.PropagationTimeout = dns2.PropagationTimeout
	cfg.PollingInterval = dns2.PollingInterval

	return luadns.NewDNSProviderConfig(cfg)
}
