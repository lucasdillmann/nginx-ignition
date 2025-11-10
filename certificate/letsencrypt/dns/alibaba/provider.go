package alibaba

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/alidns"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

const (
	accessKeyFieldID       = "alibabaAccessKeyId"
	accessKeySecretFieldID = "alibabaAccessKeySecret"
	securityTokenFieldID   = "alibabaSecurityToken"
	regionFieldID          = "alibabaRegion"
	ramRoleFieldID         = "alibabaRamRole"
)

type Provider struct{}

func (p *Provider) ID() string { return "ALIBABA" }

func (p *Provider) Name() string { return "Alibaba" }

func (p *Provider) DynamicFields() []*dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          accessKeyFieldID,
			Description: "Alibaba Cloud access key ID",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          accessKeySecretFieldID,
			Description: "Alibaba Cloud access key secret",
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          securityTokenFieldID,
			Description: "Alibaba Cloud security token",
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          regionFieldID,
			Description: "Alibaba Cloud region",
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          ramRoleFieldID,
			Description: "Alibaba Cloud RAM role",
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
	accessSecret, _ := parameters[accessKeySecretFieldID].(string)
	securityToken, _ := parameters[securityTokenFieldID].(string)
	region, _ := parameters[regionFieldID].(string)
	role, _ := parameters[ramRoleFieldID].(string)

	cfg := &alidns.Config{
		RAMRole:            role,
		APIKey:             accessKey,
		SecretKey:          accessSecret,
		SecurityToken:      securityToken,
		RegionID:           region,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PollingInterval,
		TTL:                dns.TTL,
	}

	return alidns.NewDNSProviderConfig(cfg)
}
