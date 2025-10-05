package regru

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/regru"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

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

func (p *Provider) Name() string {
	return "Reg.ru"
}

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamic_fields.DynamicField{
		{
			ID:          usernameFieldID,
			Description: "Reg.ru username",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          passwordFieldID,
			Description: "Reg.ru password",
			Required:    true,
			Sensitive:   true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          tlsCertFieldID,
			Description: "Reg.ru TLS certificate (for mTLS)",
			Sensitive:   true,
			Type:        dynamic_fields.MultiLineTextType,
		},
		{
			ID:          tlsKeyFieldID,
			Description: "Reg.ru TLS key (for mTLS)",
			Sensitive:   true,
			Type:        dynamic_fields.MultiLineTextType,
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

	cfg := &regru.Config{
		Username:           username,
		Password:           password,
		TLSCert:            tlsCert,
		TLSKey:             tlsKey,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PoolingInterval,
		TTL:                dns.TTL,
	}

	return regru.NewDNSProviderConfig(cfg)
}
