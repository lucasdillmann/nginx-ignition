package jdcloud

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/jdcloud"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

const (
	accessKeyFieldID = "jdcloudAccessKey"
	secretKeyFieldID = "jdcloudSecretKey"
	regionIDFieldID  = "jdcloudRegionId"
)

type Provider struct{}

func (p *Provider) ID() string {
	return "JD_CLOUD"
}

func (p *Provider) Name() string {
	return "JD Cloud"
}

func (p *Provider) DynamicFields() []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          accessKeyFieldID,
			Description: "JD Cloud Access Key",
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          secretKeyFieldID,
			Description: "JD Cloud Secret Key",
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          regionIDFieldID,
			Description: "JD Cloud Region ID",
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
	secretKey, _ := parameters[secretKeyFieldID].(string)
	regionID, _ := parameters[regionIDFieldID].(string)

	cfg := jdcloud.NewDefaultConfig()
	cfg.AccessKeyID = accessKey
	cfg.AccessKeySecret = secretKey
	cfg.RegionID = regionID
	cfg.TTL = dns.TTL
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval

	return jdcloud.NewDNSProviderConfig(cfg)
}
