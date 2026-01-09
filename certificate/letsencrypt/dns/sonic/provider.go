package sonic

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/sonic"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

//nolint:gosec
const (
	userIDFieldID = "sonicUserId"
	apiKeyFieldID = "sonicApiKey"
)

type Provider struct{}

func (p *Provider) ID() string {
	return "SONIC"
}

func (p *Provider) Name() string {
	return "Sonic"
}

func (p *Provider) DynamicFields() []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          userIDFieldID,
			Description: "Sonic user ID",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          apiKeyFieldID,
			Description: "Sonic API key",
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
	userID, _ := parameters[userIDFieldID].(string)
	apiKey, _ := parameters[apiKeyFieldID].(string)

	cfg := sonic.NewDefaultConfig()
	cfg.UserID = userID
	cfg.APIKey = apiKey
	cfg.TTL = dns.TTL
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval
	cfg.SequenceInterval = dns.SequenceInterval

	return sonic.NewDNSProviderConfig(cfg)
}
