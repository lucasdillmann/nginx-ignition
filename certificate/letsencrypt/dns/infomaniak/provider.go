package infomaniak

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/infomaniak"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

const (
	accessTokenFieldID = "infomaniakAccessToken"
)

type Provider struct{}

func (p *Provider) ID() string { return "INFOMANIAK" }

func (p *Provider) Name() string { return "Infomaniak" }

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamic_fields.DynamicField{
		{
			ID:          accessTokenFieldID,
			Description: "Infomaniak API access token",
			Required:    true,
			Sensitive:   true,
			Type:        dynamic_fields.SingleLineTextType,
		},
	})
}

func (p *Provider) ChallengeProvider(
	_ context.Context,
	_ []string,
	parameters map[string]any,
) (challenge.Provider, error) {
	accessToken, _ := parameters[accessTokenFieldID].(string)

	cfg := &infomaniak.Config{
		AccessToken:        accessToken,
		TTL:                dns.TTL,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PoolingInterval,
	}

	return infomaniak.NewDNSProviderConfig(cfg)
}
