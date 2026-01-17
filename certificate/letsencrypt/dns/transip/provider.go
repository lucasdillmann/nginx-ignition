package transip

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/transip"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

const (
	accountNameFieldID    = "transIpAccountName"
	privateKeyPathFieldID = "transIpPrivateKeyPath"
)

type Provider struct{}

func (p *Provider) ID() string { return "TRANSIP" }

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsTransipName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          accountNameFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsTransipAccountName),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          privateKeyPathFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsTransipPrivateKeyPath),
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
	cfg.TTL = int64(dns.TTL)
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval

	return transip.NewDNSProviderConfig(cfg)
}
