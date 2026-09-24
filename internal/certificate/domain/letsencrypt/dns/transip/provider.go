package transip

import (
	"context"

	"github.com/go-acme/lego/v5/challenge"
	"github.com/go-acme/lego/v5/providers/dns/transip"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/dynamicfields"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/i18n"
	dns2 "github.com/lucasdillmann/nginx-ignition/internal/certificate/domain/letsencrypt/dns"
)

const (
	accountNameFieldID    = "transIpAccountName"
	privateKeyPathFieldID = "transIpPrivateKeyPath"
)

type Provider struct{}

func (p *Provider) ID() string { return "TRANSIP" }

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateLetsencryptDnsTransipName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns2.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          accountNameFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsTransipAccountName),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          privateKeyPathFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsTransipPrivateKeyPath),
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
	accountName, _ := parameters[accountNameFieldID].(string)
	privateKeyPath, _ := parameters[privateKeyPathFieldID].(string)

	cfg := transip.NewDefaultConfig()
	cfg.AccountName = accountName
	cfg.PrivateKeyPath = privateKeyPath
	cfg.TTL = int64(dns2.TTL)
	cfg.PropagationTimeout = dns2.PropagationTimeout
	cfg.PollingInterval = dns2.PollingInterval

	return transip.NewDNSProviderConfig(cfg)
}
