package ovh

import (
	"context"

	"github.com/go-acme/lego/v5/challenge"
	"github.com/go-acme/lego/v5/providers/dns/ovh"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/dynamicfields"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/i18n"
	dns2 "github.com/lucasdillmann/nginx-ignition/internal/certificate/domain/letsencrypt/dns"
)

//nolint:gosec
const (
	endpointFieldID       = "ovhEndpoint"
	applicationKeyFieldID = "ovhApplicationKey"
	applicationSecFieldID = "ovhApplicationSecret"
	consumerKeyFieldID    = "ovhConsumerKey"
)

type Provider struct{}

func (p *Provider) ID() string { return "OVH" }

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateLetsencryptDnsOvhName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns2.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          endpointFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsOvhApiEndpoint),
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          applicationKeyFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsOvhApplicationKey),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          applicationSecFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsOvhApplicationSecret),
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          consumerKeyFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsOvhConsumerKey),
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
	endpoint, _ := parameters[endpointFieldID].(string)
	appKey, _ := parameters[applicationKeyFieldID].(string)
	appSecret, _ := parameters[applicationSecFieldID].(string)
	consumerKey, _ := parameters[consumerKeyFieldID].(string)

	cfg := ovh.NewDefaultConfig()
	cfg.APIEndpoint = endpoint
	cfg.ApplicationKey = appKey
	cfg.ApplicationSecret = appSecret
	cfg.ConsumerKey = consumerKey
	cfg.PropagationTimeout = dns2.PropagationTimeout
	cfg.PollingInterval = dns2.PollingInterval
	cfg.TTL = dns2.TTL

	return ovh.NewDNSProviderConfig(cfg)
}
