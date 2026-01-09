package anexia

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/anexia"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

const (
	tokenFieldID  = "anexiaToken"
	apiURLFieldID = "anexiaApiUrl"
)

type Provider struct{}

func (p *Provider) ID() string {
	return "ANEXIA"
}

func (p *Provider) Name() string {
	return "Anexia CloudDNS"
}

func (p *Provider) DynamicFields() []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          tokenFieldID,
			Description: "Anexia API token",
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          apiURLFieldID,
			Description: "Anexia API URL (optional)",
			Required:    false,
			Type:        dynamicfields.SingleLineTextType,
		},
	})
}

func (p *Provider) ChallengeProvider(
	_ context.Context,
	_ []string,
	parameters map[string]any,
) (challenge.Provider, error) {
	token, _ := parameters[tokenFieldID].(string)
	apiURL, _ := parameters[apiURLFieldID].(string)

	cfg := anexia.NewDefaultConfig()
	cfg.Token = token
	cfg.APIURL = apiURL
	cfg.TTL = dns.TTL
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval

	return anexia.NewDNSProviderConfig(cfg)
}
