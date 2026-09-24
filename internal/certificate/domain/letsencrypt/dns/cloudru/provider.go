package cloudru

import (
	"context"

	"github.com/go-acme/lego/v5/challenge"
	"github.com/go-acme/lego/v5/providers/dns/cloudru"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/dynamicfields"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/i18n"
	dns2 "github.com/lucasdillmann/nginx-ignition/internal/certificate/domain/letsencrypt/dns"
)

//nolint:gosec
const (
	serviceInstanceIDFieldID = "cloudRuServiceInstanceId"
	keyIDFieldID             = "cloudRuKeyID"
	secretFieldID            = "cloudRuSecret"
)

type Provider struct{}

func (p *Provider) ID() string { return "CLOUDRU" }

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateLetsencryptDnsCloudruName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns2.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID: serviceInstanceIDFieldID,
			Description: i18n.M(
				ctx,
				i18n.K.CertificateLetsencryptDnsCloudruServiceInstanceId,
			),
			Required: true,
			Type:     dynamicfields.SingleLineTextType,
		},
		{
			ID:          keyIDFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsCloudruKeyId),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          secretFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsCloudruSecret),
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

	cfg := cloudru.NewDefaultConfig()
	cfg.ServiceInstanceID = serviceInstanceID
	cfg.KeyID = keyID
	cfg.Secret = secret
	cfg.TTL = dns2.TTL
	cfg.PropagationTimeout = dns2.PropagationTimeout
	cfg.PollingInterval = dns2.PollingInterval
	cfg.SequenceInterval = dns2.SequenceInterval

	return cloudru.NewDNSProviderConfig(cfg)
}
