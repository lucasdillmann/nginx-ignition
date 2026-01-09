package webnamesca

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/webnamesca"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

//nolint:gosec
const (
	apiUserFieldID = "webnamescaApiUser"
	apiKeyFieldID  = "webnamescaApiKey"
)

type Provider struct{}

func (p *Provider) ID() string {
	return "WEBNAMESCA"
}

func (p *Provider) Name() string {
	return "webnames.ca"
}

func (p *Provider) DynamicFields() []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          apiUserFieldID,
			Description: "webnames.ca API user",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          apiKeyFieldID,
			Description: "webnames.ca API key",
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
	apiUser, _ := parameters[apiUserFieldID].(string)
	apiKey, _ := parameters[apiKeyFieldID].(string)

	cfg := webnamesca.NewDefaultConfig()
	cfg.APIUser = apiUser
	cfg.APIKey = apiKey
	cfg.TTL = dns.TTL
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval

	return webnamesca.NewDNSProviderConfig(cfg)
}
