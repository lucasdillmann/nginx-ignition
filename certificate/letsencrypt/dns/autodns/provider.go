package autodns

import (
	"context"
	"fmt"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/autodns"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
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
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
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
	cfg.TTL = dns.TTL
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval

	return autodns.NewDNSProviderConfig(cfg)
}
