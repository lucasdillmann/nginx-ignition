package ispconfigddns

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/ispconfigddns"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

//nolint:gosec
const (
	serverURLFieldID = "ispconfigddnsServerUrl"
	tokenFieldID     = "ispconfigddnsToken"
)

type Provider struct{}

func (p *Provider) ID() string {
	return "ISPCONFIG_DDNS"
}

func (p *Provider) Name() string {
	return "ISPConfig 3 DDNS"
}

func (p *Provider) DynamicFields() []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          serverURLFieldID,
			Description: "ISPConfig 3 DDNS server URL",
			Required:    true,
			Type:        dynamicfields.URLType,
		},
		{
			ID:          tokenFieldID,
			Description: "ISPConfig 3 DDNS token",
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
	serverURL, _ := parameters[serverURLFieldID].(string)
	token, _ := parameters[tokenFieldID].(string)

	cfg := ispconfigddns.NewDefaultConfig()
	cfg.ServerURL = serverURL
	cfg.Token = token
	cfg.TTL = dns.TTL
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval

	return ispconfigddns.NewDNSProviderConfig(cfg)
}
