package ispconfigddns

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/ispconfigddns"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

//nolint:gosec
const (
	serverURLFieldID = "ispconfigddnsServerUrl"
	tokenFieldID     = "ispconfigddnsToken"
)

type Provider struct{}

func (p *Provider) ID() string {
	return "ISPCONFIG_DDNS"
}

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateLetsencryptDnsIspconfigddnsName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          serverURLFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsIspconfigddnsServerUrl),
			Required:    true,
			Type:        dynamicfields.URLType,
		},
		{
			ID:          tokenFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsIspconfigddnsToken),
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
	serverURL, _ := parameters[serverURLFieldID].(string)
	token, _ := parameters[tokenFieldID].(string)

	cfg := ispconfigddns.NewDefaultConfig()
	cfg.ServerURL = serverURL
	cfg.Token = token
	cfg.TTL = dns.TTL
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval

	return ispconfigddns.NewDNSProviderConfig(cfg)
}
