package edgeone

import (
	"context"

	"github.com/go-acme/lego/v5/challenge"
	"github.com/go-acme/lego/v5/providers/dns/edgeone"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/dynamicfields"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/i18n"
	dns2 "github.com/lucasdillmann/nginx-ignition/internal/certificate/domain/letsencrypt/dns"
)

//nolint:gosec
const (
	secretIDFieldID  = "edgeOneSecretID"
	secretKeyFieldID = "edgeOneSecretKey"
	regionFieldID    = "edgeOneRegion"
)

type Provider struct{}

func (p *Provider) ID() string { return "EDGEONE" }

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateLetsencryptDnsEdgeoneName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns2.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          secretIDFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsEdgeoneSecretId),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          secretKeyFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsEdgeoneSecretKey),
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          regionFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsEdgeoneRegion),
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
	secretID, _ := parameters[secretIDFieldID].(string)
	secretKey, _ := parameters[secretKeyFieldID].(string)
	region, _ := parameters[regionFieldID].(string)

	cfg := edgeone.NewDefaultConfig()
	cfg.SecretID = secretID
	cfg.SecretKey = secretKey
	cfg.Region = region
	cfg.TTL = dns2.TTL
	cfg.PropagationTimeout = dns2.PropagationTimeout
	cfg.PollingInterval = dns2.PollingInterval

	return edgeone.NewDNSProviderConfig(cfg)
}
