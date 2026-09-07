package namedotcom

import (
	"context"

	"github.com/go-acme/lego/v5/challenge"
	"github.com/go-acme/lego/v5/providers/dns/namedotcom"

	"github.com/lucasdillmann/nginx-ignition/internal/certificate/acme/dns"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/dynamicfields"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/i18n"
)

const (
	usernameFieldID = "nameDotComUsername"
	apiTokenFieldID = "nameDotComApiToken"
	serverFieldID   = "nameDotComServer"
)

type Provider struct{}

func (p *Provider) ID() string { return "NAMEDOTCOM" }

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateAcmeDnsNamedotcomName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          usernameFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateAcmeDnsNamedotcomUsername),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          apiTokenFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateAcmeDnsNamedotcomApiToken),
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          serverFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateAcmeDnsNamedotcomServer),
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
	apiToken, _ := parameters[apiTokenFieldID].(string)
	server, _ := parameters[serverFieldID].(string)

	cfg := namedotcom.NewDefaultConfig()
	cfg.Username = username
	cfg.APIToken = apiToken
	cfg.Server = server
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval
	cfg.TTL = dns.TTL

	return namedotcom.NewDNSProviderConfig(cfg)
}
