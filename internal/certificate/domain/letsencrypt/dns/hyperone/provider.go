package hyperone

import (
	"context"

	"github.com/go-acme/lego/v5/challenge"
	"github.com/go-acme/lego/v5/providers/dns/hyperone"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/dynamicfields"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/i18n"
	dns2 "github.com/lucasdillmann/nginx-ignition/internal/certificate/domain/letsencrypt/dns"
)

const (
	apiEndpointFieldID      = "hyperOneApiEndpoint"
	locationIDFieldID       = "hyperOneLocationId"
	passportLocationFieldID = "hyperOnePassportLocation"
)

type Provider struct{}

func (p *Provider) ID() string { return "HYPERONE" }

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateLetsencryptDnsHyperoneName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns2.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          apiEndpointFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsHyperoneApiEndpoint),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          locationIDFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsHyperoneLocationId),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID: passportLocationFieldID,
			Description: i18n.M(
				ctx,
				i18n.K.CertificateLetsencryptDnsHyperonePassportFileLocation,
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
	cfg.TTL = dns2.TTL
	cfg.PropagationTimeout = dns2.PropagationTimeout
	cfg.PollingInterval = dns2.PollingInterval

	return hyperone.NewDNSProviderConfig(cfg)
}
