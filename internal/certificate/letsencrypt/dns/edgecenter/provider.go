package edgecenter

import (
	"context"

	"github.com/go-acme/lego/v5/challenge"
	"github.com/go-acme/lego/v5/providers/dns/edgecenter"

	"nginx-ignition/internal/certificate/letsencrypt/dns"
	"nginx-ignition/internal/core/common/dynamicfields"
	"nginx-ignition/internal/core/common/i18n"
)

const (
	tokenFieldID = "edgecenterToken"
)

type Provider struct{}

func (p *Provider) ID() string {
	return "EDGECENTER"
}

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateLetsencryptDnsEdgecenterName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          tokenFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsEdgecenterApiToken),
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
	token, _ := parameters[tokenFieldID].(string)

	cfg := edgecenter.NewDefaultConfig()
	cfg.APIToken = token
	cfg.TTL = dns.TTL
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval

	return edgecenter.NewDNSProviderConfig(cfg)
}
