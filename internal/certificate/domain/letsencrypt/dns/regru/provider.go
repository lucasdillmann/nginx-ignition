package regru

import (
	"context"

	"github.com/go-acme/lego/v5/challenge"
	"github.com/go-acme/lego/v5/providers/dns/regru"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/dynamicfields"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/i18n"
	dns2 "github.com/lucasdillmann/nginx-ignition/internal/certificate/domain/letsencrypt/dns"
)

//nolint:gosec
const (
	usernameFieldID = "regruUsername"
	passwordFieldID = "regruPassword"
	tlsCertFieldID  = "regruTlsCert"
	tlsKeyFieldID   = "regruTlsKey"
)

type Provider struct{}

func (p *Provider) ID() string {
	return "REGRU"
}

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateLetsencryptDnsRegruName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns2.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          usernameFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsRegruUsername),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          passwordFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsRegruPassword),
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          tlsCertFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsRegruTlsCertificate),
			Sensitive:   true,
			Type:        dynamicfields.MultiLineTextType,
		},
		{
			ID:          tlsKeyFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsRegruTlsKey),
			Sensitive:   true,
			Type:        dynamicfields.MultiLineTextType,
		},
	})
}

func (p *Provider) ChallengeProvider(
	_ context.Context,
	_ []string,
	parameters map[string]any,
) (challenge.Provider, error) {
	username, _ := parameters[usernameFieldID].(string)
	password, _ := parameters[passwordFieldID].(string)
	tlsCert, _ := parameters[tlsCertFieldID].(string)
	tlsKey, _ := parameters[tlsKeyFieldID].(string)

	cfg := regru.NewDefaultConfig()
	cfg.Username = username
	cfg.Password = password
	cfg.TLSCert = tlsCert
	cfg.TLSKey = tlsKey
	cfg.PropagationTimeout = dns2.PropagationTimeout
	cfg.PollingInterval = dns2.PollingInterval
	cfg.TTL = dns2.TTL

	return regru.NewDNSProviderConfig(cfg)
}
