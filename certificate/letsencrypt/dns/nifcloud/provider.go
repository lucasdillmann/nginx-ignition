package nifcloud

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/nifcloud"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

const (
	accessKeyFieldID = "nifCloudAccessKey"
	secretKeyFieldID = "nifCloudSecretKey"
)

type Provider struct{}

func (p *Provider) ID() string { return "NIFCLOUD" }

func (p *Provider) Name() string { return "NIFCLOUD" }

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamic_fields.DynamicField{
		{
			ID:          accessKeyFieldID,
			Description: "NIFCLOUD access key ID",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          secretKeyFieldID,
			Description: "NIFCLOUD secret access key",
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
	accessKey, _ := parameters[accessKeyFieldID].(string)
	secretKey, _ := parameters[secretKeyFieldID].(string)

	cfg := &nifcloud.Config{
		AccessKey:          accessKey,
		SecretKey:          secretKey,
		TTL:                dns.TTL,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PoolingInterval,
	}

	return nifcloud.NewDNSProviderConfig(cfg)
}
