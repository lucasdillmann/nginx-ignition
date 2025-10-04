package cloudru

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/cloudru"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

const (
	serviceInstanceIDFieldID = "cloudRuServiceInstanceId"
	keyIDFieldID             = "cloudRuKeyID"
	secretFieldID            = "cloudRuSecret"
)

type Provider struct{}

func (p *Provider) ID() string { return "CLOUDRU" }

func (p *Provider) Name() string { return "Cloud.ru" }

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamic_fields.DynamicField{
		{
			ID:          serviceInstanceIDFieldID,
			Description: "Cloud.ru Service Instance ID",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          keyIDFieldID,
			Description: "Cloud.ru Key ID",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          secretFieldID,
			Description: "Cloud.ru Secret",
			Required:    true,
			Sensitive:   true,
			Type:        dynamic_fields.SingleLineTextType,
		},
	})
}

func (p *Provider) ChallengeProvider(_ context.Context, _ []string, parameters map[string]any) (challenge.Provider, error) {
	serviceInstanceID, _ := parameters[serviceInstanceIDFieldID].(string)
	keyID, _ := parameters[keyIDFieldID].(string)
	secret, _ := parameters[secretFieldID].(string)

	cfg := &cloudru.Config{
		ServiceInstanceID:  serviceInstanceID,
		KeyID:              keyID,
		Secret:             secret,
		TTL:                dns.TTL,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PoolingInterval,
		SequenceInterval:   dns.SequenceInterval,
	}

	return cloudru.NewDNSProviderConfig(cfg)
}
