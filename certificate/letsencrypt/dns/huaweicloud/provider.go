package huaweicloud

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/huaweicloud"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

//nolint:gosec
const (
	accessKeyFieldID       = "huaweiCloudAccessKeyID"
	secretAccessKeyFieldID = "huaweiCloudSecretAccessKey"
	regionFieldID          = "huaweiCloudRegion"
)

type Provider struct{}

func (p *Provider) ID() string { return "HUAWEI_CLOUD" }

func (p *Provider) Name() string { return "Huawei Cloud" }

func (p *Provider) DynamicFields() []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          accessKeyFieldID,
			Description: "Huawei Cloud access key ID",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          secretAccessKeyFieldID,
			Description: "Huawei Cloud secret access key",
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          regionFieldID,
			Description: "Huawei Cloud region",
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
	accessKey, _ := parameters[accessKeyFieldID].(string)
	secretAccessKey, _ := parameters[secretAccessKeyFieldID].(string)
	region, _ := parameters[regionFieldID].(string)

	cfg := huaweicloud.NewDefaultConfig()
	cfg.AccessKeyID = accessKey
	cfg.SecretAccessKey = secretAccessKey
	cfg.Region = region
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval
	cfg.TTL = dns.TTL

	return huaweicloud.NewDNSProviderConfig(cfg)
}
