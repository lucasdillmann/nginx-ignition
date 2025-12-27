package selectel

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/selectel"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

const (
	baseURLFieldID = "selectelBaseUrl"
	tokenFieldID   = "selectelToken"
)

type Provider struct{}

func (p *Provider) ID() string {
	return "SELECTEL"
}

func (p *Provider) Name() string {
	return "Selectel"
}

func (p *Provider) DynamicFields() []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          tokenFieldID,
			Description: "Selectel API token",
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          baseURLFieldID,
			Description: "Selectel base URL",
			Type:        dynamicfields.SingleLineTextType,
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

	cfg := &selectel.Config{
		Token:              token,
		BaseURL:            baseURL,
		TTL:                dns.TTL,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PollingInterval,
	}

	return selectel.NewDNSProviderConfig(cfg)
}
