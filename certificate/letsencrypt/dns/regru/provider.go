package regru

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/regru"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
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

func (p *Provider) DynamicFields() []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          usernameFieldID,
			Description: "Reg.ru username",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          passwordFieldID,
			Description: "Reg.ru password",
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          tlsCertFieldID,
			Description: "Reg.ru TLS certificate (for mTLS)",
			Sensitive:   true,
			Type:        dynamicfields.MultiLineTextType,
		},
		{
			ID:          tlsKeyFieldID,
			Description: "Reg.ru TLS key (for mTLS)",
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

	cfg := &regru.Config{
		Username:           username,
		Password:           password,
		TLSCert:            tlsCert,
		TLSKey:             tlsKey,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PollingInterval,
		TTL:                dns.TTL,
	}

	return regru.NewDNSProviderConfig(cfg)
}
