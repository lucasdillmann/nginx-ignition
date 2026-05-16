package gehirn

import (
	"context"

	"github.com/go-acme/lego/v5/challenge"
	"github.com/go-acme/lego/v5/providers/dns/gehirn"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

const (
	tokenIDFieldID     = "gehirnTokenId"     // nolint:gosec
	tokenSecretFieldID = "gehirnTokenSecret" // nolint:gosec
)

type Provider struct{}

func (p *Provider) ID() string {
	return "GEHIRN"
}

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateLetsencryptDnsGehirnName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          tokenIDFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsGehirnTokenId),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          tokenSecretFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsGehirnTokenSecret),
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
	tokenID, _ := parameters[tokenIDFieldID].(string)
	tokenSecret, _ := parameters[tokenSecretFieldID].(string)

	cfg := gehirn.NewDefaultConfig()
	cfg.TokenID = tokenID
	cfg.TokenSecret = tokenSecret
	cfg.TTL = dns.TTL
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval

	return gehirn.NewDNSProviderConfig(cfg)
}
