package easydns

import (
	"context"
	"net/url"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/easydns"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
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
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
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
	cfg.TTL = dns.TTL
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.SequenceInterval = dns.SequenceInterval
	cfg.PollingInterval = dns.PollingInterval

	if endpoint != "" {
		if parsedValue, err := url.Parse(endpoint); err == nil {
			cfg.Endpoint = parsedValue
		}
	}

	return easydns.NewDNSProviderConfig(cfg)
}
