package yandex

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/yandex"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

const (
	pddTokenFieldID = "yandexPddToken"
)

type Provider struct{}

func (p *Provider) ID() string { return "YANDEX" }

func (p *Provider) Name() string { return "Yandex" }

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamic_fields.DynamicField{
		{
			ID:          pddTokenFieldID,
			Description: "Yandex PDD token",
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
	pddToken, _ := parameters[pddTokenFieldID].(string)

	cfg := &yandex.Config{
		PddToken:           pddToken,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PoolingInterval,
		TTL:                dns.TTL,
	}

	return yandex.NewDNSProviderConfig(cfg)
}
