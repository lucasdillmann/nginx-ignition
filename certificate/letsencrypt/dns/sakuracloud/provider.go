package sakuracloud

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/sakuracloud"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

const (
	accessTokenFieldID  = "sakuraCloudAccessToken"
	accessSecretFieldID = "sakuraCloudAccessSecret"
)

type Provider struct{}

func (p *Provider) ID() string { return "SAKURA_CLOUD" }

func (p *Provider) Name() string { return "SakuraCloud" }

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamic_fields.DynamicField{
		{
			ID:          accessTokenFieldID,
			Description: "SakuraCloud access token",
			Required:    true,
			Sensitive:   true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          accessSecretFieldID,
			Description: "SakuraCloud access secret",
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
	accessToken, _ := parameters[accessTokenFieldID].(string)
	accessSecret, _ := parameters[accessSecretFieldID].(string)

	cfg := &sakuracloud.Config{
		Token:              accessToken,
		Secret:             accessSecret,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PoolingInterval,
		TTL:                dns.TTL,
	}

	return sakuracloud.NewDNSProviderConfig(cfg)
}
