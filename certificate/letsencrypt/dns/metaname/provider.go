package metaname

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/metaname"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

//nolint:gosec
const (
	accountReferenceFieldID = "metanameAccountReference"
	apiKeyFieldID           = "metanameApiKey"
)

type Provider struct{}

func (p *Provider) ID() string { return "METANAME" }

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsMetanameName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID: accountReferenceFieldID,
			Description: i18n.M(
				ctx,
				i18n.K.CertificateCommonLetsEncryptDnsMetanameAccountReference,
			),
			Required: true,
			Type:     dynamicfields.SingleLineTextType,
		},
		{
			ID:          apiKeyFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsMetanameApiKey),
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
	accountReference, _ := parameters[accountReferenceFieldID].(string)
	apiKey, _ := parameters[apiKeyFieldID].(string)

	cfg := metaname.NewDefaultConfig()
	cfg.AccountReference = accountReference
	cfg.APIKey = apiKey
	cfg.TTL = dns.TTL
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval

	return metaname.NewDNSProviderConfig(cfg)
}
