package netcup

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/netcup"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

//nolint:gosec
const (
	customerFieldID    = "netcupCustomer"
	apiKeyFieldID      = "netcupApiKey"
	apiPasswordFieldID = "netcupApiPassword"
)

type Provider struct{}

func (p *Provider) ID() string { return "NETCUP" }

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateLetsencryptDnsNetcupName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          customerFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsNetcupCustomer),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          apiKeyFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsNetcupApiKey),
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          apiPasswordFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsNetcupApiPassword),
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
	customerNumber, _ := parameters[customerFieldID].(string)
	apiKey, _ := parameters[apiKeyFieldID].(string)
	apiPassword, _ := parameters[apiPasswordFieldID].(string)

	cfg := netcup.NewDefaultConfig()
	cfg.Customer = customerNumber
	cfg.Key = apiKey
	cfg.Password = apiPassword
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval

	return netcup.NewDNSProviderConfig(cfg)
}
