package rcodezero

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/rcodezero"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

//nolint:gosec
const (
	apiTokenFieldID = "rcodeZeroApiToken"
)

type Provider struct{}

func (p *Provider) ID() string {
	return "RCODEZERO"
}

func (p *Provider) Name() string {
	return "RcodeZero"
}

func (p *Provider) DynamicFields() []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          apiTokenFieldID,
			Description: "RcodeZero API token",
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
	apiToken, _ := parameters[apiTokenFieldID].(string)

	cfg := &rcodezero.Config{
		APIToken:           apiToken,
		TTL:                dns.TTL,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PollingInterval,
	}

	return rcodezero.NewDNSProviderConfig(cfg)
}
