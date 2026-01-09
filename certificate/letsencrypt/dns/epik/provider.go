package epik

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/epik"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

const (
	signatureFieldID = "epikSignature"
)

type Provider struct{}

func (p *Provider) ID() string { return "EPIK" }

func (p *Provider) Name() string { return "Epik" }

func (p *Provider) DynamicFields() []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          signatureFieldID,
			Description: "Epik API signature",
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
	signature, _ := parameters[signatureFieldID].(string)

	cfg := epik.NewDefaultConfig()
	cfg.Signature = signature
	cfg.TTL = dns.TTL
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval

	return epik.NewDNSProviderConfig(cfg)
}
