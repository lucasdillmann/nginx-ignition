package mydnsjp

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/mydnsjp"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

const (
	masterIDFieldID = "myDnsJpMasterId"
	passwordFieldID = "myDnsJpPassword"
)

type Provider struct{}

func (p *Provider) ID() string { return "MYDNS_JP" }

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateLetsencryptDnsMydnsjpName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          masterIDFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsMydnsjpMasterId),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          passwordFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsMydnsjpPassword),
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
	masterID, _ := parameters[masterIDFieldID].(string)
	password, _ := parameters[passwordFieldID].(string)

	cfg := mydnsjp.NewDefaultConfig()
	cfg.MasterID = masterID
	cfg.Password = password
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval

	return mydnsjp.NewDNSProviderConfig(cfg)
}
