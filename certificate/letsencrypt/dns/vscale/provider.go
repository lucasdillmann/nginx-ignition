package vscale

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/vscale"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

const (
	baseURLFieldID = "vscaleBaseUrl"
	tokenFieldID   = "vscaleToken"
)

type Provider struct{}

func (p *Provider) ID() string {
	return "VSCALE"
}

func (p *Provider) Name() string {
	return "Vscale"
}

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamic_fields.DynamicField{
		{
			ID:          tokenFieldID,
			Description: "Vscale API token",
			Required:    true,
			Sensitive:   true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          baseURLFieldID,
			Description: "Vscale base URL",
			Type:        dynamic_fields.SingleLineTextType,
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

	cfg := &vscale.Config{
		Token:              token,
		BaseURL:            baseURL,
		TTL:                dns.TTL,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PoolingInterval,
	}

	return vscale.NewDNSProviderConfig(cfg)
}
