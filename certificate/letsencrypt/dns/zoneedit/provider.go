package zoneedit

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/zoneedit"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

//nolint:gosec
const (
	userFieldID      = "zoneeditUser"
	authTokenFieldID = "zoneeditAuthToken"
)

type Provider struct{}

func (p *Provider) ID() string {
	return "ZONEEDIT"
}

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateLetsencryptDnsZoneeditName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          userFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsZoneeditUser),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          authTokenFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsZoneeditAuthToken),
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
	user, _ := parameters[userFieldID].(string)
	authToken, _ := parameters[authTokenFieldID].(string)

	cfg := zoneedit.NewDefaultConfig()
	cfg.User = user
	cfg.AuthToken = authToken
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval

	return zoneedit.NewDNSProviderConfig(cfg)
}
