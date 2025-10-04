package baiducloud

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/baiducloud"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

const (
	accessKeyFieldID       = "baiduCloudAccessKeyID"
	secretAccessKeyFieldID = "baiduCloudSecretAccessKey"
)

type Provider struct{}

func (p *Provider) ID() string { return "BAIDU_CLOUD" }

func (p *Provider) Name() string { return "Baidu Cloud" }

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamic_fields.DynamicField{
		{
			ID:          accessKeyFieldID,
			Description: "Baidu Cloud access key ID",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          secretAccessKeyFieldID,
			Description: "Baidu Cloud secret access key",
			Required:    true,
			Sensitive:   true,
			Type:        dynamic_fields.SingleLineTextType,
		},
	})
}

func (p *Provider) ChallengeProvider(_ context.Context, _ []string, parameters map[string]any) (challenge.Provider, error) {
	accessKey, _ := parameters[accessKeyFieldID].(string)
	secretKey, _ := parameters[secretAccessKeyFieldID].(string)

	cfg := &baiducloud.Config{
		AccessKeyID:        accessKey,
		SecretAccessKey:    secretKey,
		TTL:                dns.TTL,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PoolingInterval,
	}

	return baiducloud.NewDNSProviderConfig(cfg)
}
