package euserv

import (
	"context"

	"github.com/go-acme/lego/v5/challenge"
	"github.com/go-acme/lego/v5/providers/dns/euserv"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

const (
	emailFieldID    = "euservEmail"
	passwordFieldID = "euservPassword"
	orderIDFieldID  = "euservOrderId"
)

type Provider struct{}

func (p *Provider) ID() string {
	return "EUSERV"
}

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateLetsencryptDnsEuservName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          emailFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsEuservEmail),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          passwordFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsEuservPassword),
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          orderIDFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsEuservOrderId),
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
	email, _ := parameters[emailFieldID].(string)
	password, _ := parameters[passwordFieldID].(string)
	orderID, _ := parameters[orderIDFieldID].(string)

	cfg := euserv.NewDefaultConfig()
	cfg.Email = email
	cfg.Password = password
	cfg.OrderID = orderID
	cfg.TTL = dns.TTL
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval

	return euserv.NewDNSProviderConfig(cfg)
}
