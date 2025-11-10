package metaname

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/metaname"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

const (
	accountReferenceFieldID = "metanameAccountReference"
	apiKeyFieldID           = "metanameApiKey"
)

type Provider struct{}

func (p *Provider) ID() string { return "METANAME" }

func (p *Provider) Name() string { return "Metaname" }

func (p *Provider) DynamicFields() []*dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          accountReferenceFieldID,
			Description: "Metaname account reference",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          apiKeyFieldID,
			Description: "Metaname API key",
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
	accountReference, _ := parameters[accountReferenceFieldID].(string)
	apiKey, _ := parameters[apiKeyFieldID].(string)

	cfg := &metaname.Config{
		AccountReference:   accountReference,
		APIKey:             apiKey,
		TTL:                dns.TTL,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PollingInterval,
	}

	return metaname.NewDNSProviderConfig(cfg)
}
