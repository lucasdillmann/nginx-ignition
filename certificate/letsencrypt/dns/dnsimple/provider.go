package dnsimple

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/dnsimple"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

const (
	accessTokenFieldID = "dnSimpleAccessToken"
	baseURLFieldID     = "dnSimpleBaseUrl"
)

type Provider struct{}

func (p *Provider) ID() string { return "DNSIMPLE" }

func (p *Provider) Name() string { return "DNSimple" }

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamic_fields.DynamicField{
		{
			ID:          accessTokenFieldID,
			Description: "DNSimple access token",
			Required:    true,
			Sensitive:   true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          baseURLFieldID,
			Description: "DNSimple base URL",
			Type:        dynamic_fields.SingleLineTextType,
		},
	})
}

func (p *Provider) ChallengeProvider(_ context.Context, _ []string, parameters map[string]any) (challenge.Provider, error) {
	accessToken, _ := parameters[accessTokenFieldID].(string)
	baseURL, _ := parameters[baseURLFieldID].(string)

	cfg := &dnsimple.Config{
		AccessToken:        accessToken,
		BaseURL:            baseURL,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PoolingInterval,
		TTL:                dns.TTL,
	}

	return dnsimple.NewDNSProviderConfig(cfg)
}
