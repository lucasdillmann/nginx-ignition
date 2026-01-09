package nearlyfreespeech

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/nearlyfreespeech"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

//nolint:gosec
const (
	loginFieldID  = "nearlyFreeSpeechLogin"
	apiKeyFieldID = "nearlyFreeSpeechApiKey"
)

type Provider struct{}

func (p *Provider) ID() string { return "NEARLYFREESPEECH" }

func (p *Provider) Name() string { return "NearlyFreeSpeech.NET" }

func (p *Provider) DynamicFields() []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          loginFieldID,
			Description: "NearlyFreeSpeech.NET login",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          apiKeyFieldID,
			Description: "NearlyFreeSpeech.NET API key",
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
	login, _ := parameters[loginFieldID].(string)
	apiKey, _ := parameters[apiKeyFieldID].(string)

	cfg := nearlyfreespeech.NewDefaultConfig()
	cfg.Login = login
	cfg.APIKey = apiKey
	cfg.TTL = dns.TTL
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval
	cfg.SequenceInterval = dns.SequenceInterval

	return nearlyfreespeech.NewDNSProviderConfig(cfg)
}
