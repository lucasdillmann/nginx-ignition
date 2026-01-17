package iij

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/iij"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

//nolint:gosec
const (
	accessKeyFieldID     = "iijAccessKey"
	secretKeyFieldID     = "iijSecretKey"
	doServiceCodeFieldID = "iijDoServiceCode"
)

type Provider struct{}

func (p *Provider) ID() string { return "IIJ" }

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsIijName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          accessKeyFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsIijAccessKey),
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          secretKeyFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsIijSecretKey),
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          doServiceCodeFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsIijServiceCode),
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
	accessKey, _ := parameters[accessKeyFieldID].(string)
	secretKey, _ := parameters[secretKeyFieldID].(string)
	doServiceCode, _ := parameters[doServiceCodeFieldID].(string)

	cfg := iij.NewDefaultConfig()
	cfg.AccessKey = accessKey
	cfg.SecretKey = secretKey
	cfg.DoServiceCode = doServiceCode
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval
	cfg.TTL = dns.TTL

	return iij.NewDNSProviderConfig(cfg)
}
