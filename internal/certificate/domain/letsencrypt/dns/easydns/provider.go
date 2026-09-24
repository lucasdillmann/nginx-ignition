package easydns

import (
	"context"
	"net/url"

	"github.com/go-acme/lego/v5/challenge"
	"github.com/go-acme/lego/v5/providers/dns/easydns"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/dynamicfields"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/i18n"
	dns2 "github.com/lucasdillmann/nginx-ignition/internal/certificate/domain/letsencrypt/dns"
)

//nolint:gosec
const (
	endpointFieldID = "easyDnsEndpoint"
	tokenFieldID    = "easyDnsToken"
	keyFieldID      = "easyDnsKey"
)

type Provider struct{}

func (p *Provider) ID() string { return "EASYDNS" }

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateLetsencryptDnsEasydnsName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns2.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          tokenFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsEasydnsApiToken),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          keyFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsEasydnsApiKey),
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          endpointFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsEasydnsApiEndpoint),
			Type:        dynamicfields.SingleLineTextType,
		},
	})
}

func (p *Provider) ChallengeProvider(
	_ context.Context,
	_ []string,
	parameters map[string]any,
) (challenge.Provider, error) {
	token, _ := parameters[tokenFieldID].(string)
	apiKey, _ := parameters[keyFieldID].(string)
	endpoint, _ := parameters[endpointFieldID].(string)

	cfg := easydns.NewDefaultConfig()
	cfg.Token = token
	cfg.Key = apiKey
	cfg.TTL = dns2.TTL
	cfg.PropagationTimeout = dns2.PropagationTimeout
	cfg.SequenceInterval = dns2.SequenceInterval
	cfg.PollingInterval = dns2.PollingInterval

	if endpoint != "" {
		if parsedValue, err := url.Parse(endpoint); err == nil {
			cfg.Endpoint = parsedValue
		}
	}

	return easydns.NewDNSProviderConfig(cfg)
}
