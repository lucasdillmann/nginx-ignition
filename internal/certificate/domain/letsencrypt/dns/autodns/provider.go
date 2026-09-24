package autodns

import (
	"context"
	"fmt"

	"github.com/go-acme/lego/v5/challenge"
	"github.com/go-acme/lego/v5/providers/dns/autodns"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/dynamicfields"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/i18n"
	dns2 "github.com/lucasdillmann/nginx-ignition/internal/certificate/domain/letsencrypt/dns"
)

const (
	usernameFieldID = "autoDnsUsername"
	passwordFieldID = "autoDnsPassword"
	contextFieldID  = "autoDnsContext"
)

type Provider struct{}

func (p *Provider) ID() string { return "AUTODNS" }

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateLetsencryptDnsAutodnsName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns2.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          usernameFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsAutodnsUsername),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          passwordFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsAutodnsPassword),
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          contextFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsAutodnsContext),
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
	username, _ := parameters[usernameFieldID].(string)
	password, _ := parameters[passwordFieldID].(string)
	contextStr, _ := parameters[contextFieldID].(string)

	var contextInt int
	if contextStr != "" {
		_, err := fmt.Sscanf(contextStr, "%d", &contextInt)
		if err != nil {
			return nil, err
		}
	}

	cfg := autodns.NewDefaultConfig()
	cfg.Username = username
	cfg.Password = password
	cfg.Context = contextInt
	cfg.TTL = dns2.TTL
	cfg.PropagationTimeout = dns2.PropagationTimeout
	cfg.PollingInterval = dns2.PollingInterval

	return autodns.NewDNSProviderConfig(cfg)
}
