package mydnsjp

import (
	"context"

	"github.com/go-acme/lego/v5/challenge"
	"github.com/go-acme/lego/v5/providers/dns/mydnsjp"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/dynamicfields"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/i18n"
	dns2 "github.com/lucasdillmann/nginx-ignition/internal/certificate/domain/letsencrypt/dns"
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
	return dns2.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
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
	cfg.PropagationTimeout = dns2.PropagationTimeout
	cfg.PollingInterval = dns2.PollingInterval

	return mydnsjp.NewDNSProviderConfig(cfg)
}
