package xinnet

import (
	"context"

	"github.com/go-acme/lego/v5/challenge"
	"github.com/go-acme/lego/v5/providers/dns/xinnet"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

const (
	secretFieldID  = "xinnetSecret"
	agentIDFieldID = "xinnetAgentId"
)

type Provider struct{}

func (p *Provider) ID() string {
	return "XINNET"
}

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateLetsencryptDnsXinnetName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          secretFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsXinnetSecret),
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          agentIDFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsXinnetAgentId),
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
	secret, _ := parameters[secretFieldID].(string)
	agentID, _ := parameters[agentIDFieldID].(string)

	cfg := xinnet.NewDefaultConfig()
	cfg.Secret = secret
	cfg.AgentID = agentID
	cfg.TTL = dns.TTL
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval

	return xinnet.NewDNSProviderConfig(cfg)
}
