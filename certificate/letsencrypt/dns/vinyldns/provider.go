package vinyldns

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/vinyldns"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

//nolint:gosec
const (
	accessKeyFieldID = "vinylDnsAccessKey"
	secretKeyFieldID = "vinylDnsSecretKey"
	hostFieldID      = "vinylDnsHost"
)

type Provider struct{}

func (p *Provider) ID() string {
	return "VINYLDNS"
}

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsVinyldnsName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          accessKeyFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsVinyldnsAccessKey),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          secretKeyFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsVinyldnsSecretKey),
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          hostFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsVinyldnsHost),
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
	host, _ := parameters[hostFieldID].(string)

	cfg := vinyldns.NewDefaultConfig()
	cfg.AccessKey = accessKey
	cfg.SecretKey = secretKey
	cfg.Host = host
	cfg.TTL = dns.TTL
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval

	return vinyldns.NewDNSProviderConfig(cfg)
}
