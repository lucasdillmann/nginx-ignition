package manageengine

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/manageengine"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

//nolint:gosec
const (
	clientIDFieldID     = "manageEngineClientId"
	clientSecretFieldID = "manageEngineClientSecret"
)

type Provider struct{}

func (p *Provider) ID() string { return "MANAGEENGINE" }

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsManageengineName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          clientIDFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsManageengineClientId),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID: clientSecretFieldID,
			Description: i18n.M(
				ctx,
				i18n.K.CertificateCommonLetsEncryptDnsManageengineClientSecret,
			),
			Required:  true,
			Sensitive: true,
			Type:      dynamicfields.SingleLineTextType,
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

	cfg := manageengine.NewDefaultConfig()
	cfg.ClientID = clientID
	cfg.ClientSecret = clientSecret
	cfg.TTL = dns.TTL
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval

	return manageengine.NewDNSProviderConfig(cfg)
}
