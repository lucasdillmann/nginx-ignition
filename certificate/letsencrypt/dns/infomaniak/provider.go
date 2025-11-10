package infomaniak

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/infomaniak"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

const (
	accessTokenFieldID = "infomaniakAccessToken"
)

type Provider struct{}

func (p *Provider) ID() string { return "INFOMANIAK" }

func (p *Provider) Name() string { return "Infomaniak" }

func (p *Provider) DynamicFields() []*dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          accessTokenFieldID,
			Description: "Infomaniak API access token",
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
	accessToken, _ := parameters[accessTokenFieldID].(string)

	cfg := &infomaniak.Config{
		AccessToken:        accessToken,
		TTL:                dns.TTL,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PollingInterval,
	}

	return infomaniak.NewDNSProviderConfig(cfg)
}
