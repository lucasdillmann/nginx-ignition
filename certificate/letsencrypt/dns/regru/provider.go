package regru

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/regru"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
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
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
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
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval
	cfg.TTL = dns.TTL

	return regru.NewDNSProviderConfig(cfg)
}
