package bindman

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/bindman"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

const (
	baseURLFieldID = "bindmanBaseURL"
)

type Provider struct{}

func (p *Provider) ID() string { return "BINDMAN" }

func (p *Provider) Name() string { return "Bindman" }

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamic_fields.DynamicField{
		{
			ID:          baseURLFieldID,
			Description: "Bindman API base URL",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
		},
	})
}

func (p *Provider) ChallengeProvider(
	_ context.Context,
	_ []string,
	parameters map[string]any,
) (challenge.Provider, error) {
	baseURL, _ := parameters[baseURLFieldID].(string)

	cfg := &bindman.Config{
		BaseURL:            baseURL,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PoolingInterval,
	}

	return bindman.NewDNSProviderConfig(cfg)
}
