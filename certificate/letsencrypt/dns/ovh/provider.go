package ovh

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/ovh"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
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
	return i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsOvhName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          endpointFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsOvhApiEndpoint),
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          applicationKeyFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsOvhApplicationKey),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          applicationSecFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsOvhApplicationSecret),
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          consumerKeyFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsOvhConsumerKey),
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
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval
	cfg.TTL = dns.TTL

	return ovh.NewDNSProviderConfig(cfg)
}
