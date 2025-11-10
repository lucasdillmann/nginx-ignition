package cloudru

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/cloudru"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

const (
	serviceInstanceIDFieldID = "cloudRuServiceInstanceId"
	keyIDFieldID             = "cloudRuKeyID"
	secretFieldID            = "cloudRuSecret"
)

type Provider struct{}

func (p *Provider) ID() string { return "CLOUDRU" }

func (p *Provider) Name() string { return "Cloud.ru" }

func (p *Provider) DynamicFields() []*dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          serviceInstanceIDFieldID,
			Description: "Cloud.ru service instance ID",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          keyIDFieldID,
			Description: "Cloud.ru key ID",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          secretFieldID,
			Description: "Cloud.ru secret",
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
	serviceInstanceID, _ := parameters[serviceInstanceIDFieldID].(string)
	keyID, _ := parameters[keyIDFieldID].(string)
	secret, _ := parameters[secretFieldID].(string)

	cfg := &cloudru.Config{
		ServiceInstanceID:  serviceInstanceID,
		KeyID:              keyID,
		Secret:             secret,
		TTL:                dns.TTL,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PollingInterval,
		SequenceInterval:   dns.SequenceInterval,
	}

	return cloudru.NewDNSProviderConfig(cfg)
}
