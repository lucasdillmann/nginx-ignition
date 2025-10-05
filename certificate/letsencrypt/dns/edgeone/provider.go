package edgeone

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/edgeone"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

const (
	secretIDFieldID  = "edgeOneSecretID"
	secretKeyFieldID = "edgeOneSecretKey"
	regionFieldID    = "edgeOneRegion"
)

type Provider struct{}

func (p *Provider) ID() string { return "EDGEONE" }

func (p *Provider) Name() string { return "Tencent EdgeOne" }

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamic_fields.DynamicField{
		{
			ID:          secretIDFieldID,
			Description: "Tencent Cloud secret ID",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          secretKeyFieldID,
			Description: "Tencent Cloud secret key",
			Required:    true,
			Sensitive:   true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          regionFieldID,
			Description: "Tencent Cloud region",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
		},
	})
}

func (p *Provider) ChallengeProvider(_ context.Context, _ []string, parameters map[string]any) (challenge.Provider, error) {
	secretID, _ := parameters[secretIDFieldID].(string)
	secretKey, _ := parameters[secretKeyFieldID].(string)
	region, _ := parameters[regionFieldID].(string)

	cfg := &edgeone.Config{
		SecretID:           secretID,
		SecretKey:          secretKey,
		Region:             region,
		TTL:                dns.TTL,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PoolingInterval,
	}

	return edgeone.NewDNSProviderConfig(cfg)
}
