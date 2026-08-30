package timewebcloud

import (
	"context"

	"github.com/go-acme/lego/v5/challenge"
	"github.com/go-acme/lego/v5/providers/dns/timewebcloud"

	"github.com/lucasdillmann/nginx-ignition/internal/certificate/letsencrypt/dns"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/dynamicfields"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/i18n"
)

//nolint:gosec
const (
	authTokenFieldID = "timewebCloudAuthToken"
)

type Provider struct{}

func (p *Provider) ID() string {
	return "TIMEWEBCLOUD"
}

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateLetsencryptDnsTimewebcloudName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          authTokenFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsTimewebcloudAuthToken),
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
	authToken, _ := parameters[authTokenFieldID].(string)

	cfg := timewebcloud.NewDefaultConfig()
	cfg.AuthToken = authToken
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval

	return timewebcloud.NewDNSProviderConfig(cfg)
}
