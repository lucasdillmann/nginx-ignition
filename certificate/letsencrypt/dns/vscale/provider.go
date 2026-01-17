package vscale

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/vscale"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

const (
	baseURLFieldID = "vscaleBaseUrl"
	tokenFieldID   = "vscaleToken"
)

type Provider struct{}

func (p *Provider) ID() string {
	return "VSCALE"
}

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsVscaleName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          tokenFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsVscaleToken),
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          baseURLFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsVscaleBaseUrl),
			Type:        dynamicfields.SingleLineTextType,
		},
	})
}

func (p *Provider) ChallengeProvider(
	_ context.Context,
	_ []string,
	parameters map[string]any,
) (challenge.Provider, error) {
	token, _ := parameters[tokenFieldID].(string)
	baseURL, _ := parameters[baseURLFieldID].(string)

	cfg := vscale.NewDefaultConfig()
	cfg.Token = token
	cfg.BaseURL = baseURL
	cfg.TTL = dns.TTL
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval

	return vscale.NewDNSProviderConfig(cfg)
}
