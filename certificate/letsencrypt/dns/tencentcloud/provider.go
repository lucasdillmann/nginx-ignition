package tencentcloud

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/tencentcloud"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

//nolint:gosec
const (
	secretIDKeyFieldID  = "tencentCloudSecretId"
	secretKeyKeyFieldID = "tencentCloudSecretKey"
	regionFieldID       = "tencentCloudRegion"
)

type Provider struct{}

func (p *Provider) ID() string { return "TENCENT_CLOUD" }

func (p *Provider) Name() string { return "Tencent Cloud" }

func (p *Provider) DynamicFields() []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          secretIDKeyFieldID,
			Description: "Tencent Cloud secret ID",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          secretKeyKeyFieldID,
			Description: "Tencent Cloud secret key",
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          regionFieldID,
			Description: "Tencent Cloud region",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
	})
}

func (p *Provider) ChallengeProvider(
	_ context.Context,
	_ []string,
	parameters map[string]any,
) (challenge.Provider, error) {
	secretID, _ := parameters[secretIDKeyFieldID].(string)
	secretKey, _ := parameters[secretKeyKeyFieldID].(string)
	region, _ := parameters[regionFieldID].(string)

	cfg := tencentcloud.NewDefaultConfig()
	cfg.SecretID = secretID
	cfg.SecretKey = secretKey
	cfg.Region = region
	cfg.TTL = dns.TTL
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval

	return tencentcloud.NewDNSProviderConfig(cfg)
}
