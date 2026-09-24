package tele3

import (
	"context"

	"github.com/go-acme/lego/v5/challenge"
	"github.com/go-acme/lego/v5/providers/dns/tele3"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/dynamicfields"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/i18n"
	dns2 "github.com/lucasdillmann/nginx-ignition/internal/certificate/domain/letsencrypt/dns"
)

const (
	keyFieldID    = "tele3Key"
	secretFieldID = "tele3Secret"
)

type Provider struct{}

func (p *Provider) ID() string {
	return "TELE3"
}

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateLetsencryptDnsTele3Name)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns2.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          keyFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsTele3Key),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          secretFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsTele3Secret),
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
	key, _ := parameters[keyFieldID].(string)
	secret, _ := parameters[secretFieldID].(string)

	cfg := tele3.NewDefaultConfig()
	cfg.Key = key
	cfg.Secret = secret
	cfg.TTL = dns2.TTL
	cfg.PropagationTimeout = dns2.PropagationTimeout
	cfg.PollingInterval = dns2.PollingInterval

	return tele3.NewDNSProviderConfig(cfg)
}
