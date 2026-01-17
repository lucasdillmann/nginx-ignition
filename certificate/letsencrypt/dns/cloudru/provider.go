package cloudru

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/cloudru"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
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
	return i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsCloudruName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID: serviceInstanceIDFieldID,
			Description: i18n.M(
				ctx,
				i18n.K.CertificateCommonLetsEncryptDnsCloudruServiceInstanceId,
			),
			Required: true,
			Type:     dynamicfields.SingleLineTextType,
		},
		{
			ID:          keyIDFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsCloudruKeyId),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          secretFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsCloudruSecret),
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
	cfg.TTL = dns.TTL
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval
	cfg.SequenceInterval = dns.SequenceInterval

	return cloudru.NewDNSProviderConfig(cfg)
}
