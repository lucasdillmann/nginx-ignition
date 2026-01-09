package yandex

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/yandex"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

//nolint:gosec
const (
	pddTokenFieldID = "yandexPddToken"
)

type Provider struct{}

func (p *Provider) ID() string { return "YANDEX" }

func (p *Provider) Name() string { return "Yandex" }

func (p *Provider) DynamicFields() []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          pddTokenFieldID,
			Description: "Yandex PDD token",
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
	pddToken, _ := parameters[pddTokenFieldID].(string)

	cfg := yandex.NewDefaultConfig()
	cfg.PddToken = pddToken
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval
	cfg.TTL = dns.TTL

	return yandex.NewDNSProviderConfig(cfg)
}
