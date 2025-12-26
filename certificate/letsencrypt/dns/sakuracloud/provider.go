package sakuracloud

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/sakuracloud"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

const (
	accessTokenFieldID  = "sakuraCloudAccessToken"
	accessSecretFieldID = "sakuraCloudAccessSecret"
)

type Provider struct{}

func (p *Provider) ID() string { return "SAKURA_CLOUD" }

func (p *Provider) Name() string { return "SakuraCloud" }

func (p *Provider) DynamicFields() []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          accessTokenFieldID,
			Description: "SakuraCloud access token",
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          accessSecretFieldID,
			Description: "SakuraCloud access secret",
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
	accessSecret, _ := parameters[accessSecretFieldID].(string)

	cfg := &sakuracloud.Config{
		Token:              accessToken,
		Secret:             accessSecret,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PollingInterval,
		TTL:                dns.TTL,
	}

	return sakuracloud.NewDNSProviderConfig(cfg)
}
