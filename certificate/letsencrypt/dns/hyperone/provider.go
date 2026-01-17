package hyperone

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/hyperone"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

const (
	apiEndpointFieldID      = "hyperOneApiEndpoint"
	locationIDFieldID       = "hyperOneLocationId"
	passportLocationFieldID = "hyperOnePassportLocation"
)

type Provider struct{}

func (p *Provider) ID() string { return "HYPERONE" }

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsHyperoneName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          apiEndpointFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsHyperoneApiEndpoint),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          locationIDFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsHyperoneLocationId),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID: passportLocationFieldID,
			Description: i18n.M(
				ctx,
				i18n.K.CertificateCommonLetsEncryptDnsHyperonePassportFileLocation,
			),
			Required: true,
			Type:     dynamicfields.SingleLineTextType,
		},
	})
}

func (p *Provider) ChallengeProvider(
	_ context.Context,
	_ []string,
	parameters map[string]any,
) (challenge.Provider, error) {
	apiEndpoint, _ := parameters[apiEndpointFieldID].(string)
	locationID, _ := parameters[locationIDFieldID].(string)
	passportLocation, _ := parameters[passportLocationFieldID].(string)

	cfg := hyperone.NewDefaultConfig()
	cfg.APIEndpoint = apiEndpoint
	cfg.LocationID = locationID
	cfg.PassportLocation = passportLocation
	cfg.TTL = dns.TTL
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval

	return hyperone.NewDNSProviderConfig(cfg)
}
