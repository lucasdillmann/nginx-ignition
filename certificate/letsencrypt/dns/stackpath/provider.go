package stackpath

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/stackpath"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

//nolint:gosec
const (
	clientIDFieldID     = "stackPathClientId"
	clientSecretFieldID = "stackPathClientSecret"
	stackIDFieldID      = "stackPathStackId"
)

type Provider struct{}

func (p *Provider) ID() string {
	return "STACKPATH"
}

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsStackpathName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          clientIDFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsStackpathClientId),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          clientSecretFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsStackpathClientSecret),
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          stackIDFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsStackpathStackId),
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
	clientID, _ := parameters[clientIDFieldID].(string)
	clientSecret, _ := parameters[clientSecretFieldID].(string)
	stackID, _ := parameters[stackIDFieldID].(string)

	cfg := stackpath.NewDefaultConfig()
	cfg.ClientID = clientID
	cfg.ClientSecret = clientSecret
	cfg.StackID = stackID
	cfg.TTL = dns.TTL
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval

	return stackpath.NewDNSProviderConfig(cfg)
}
