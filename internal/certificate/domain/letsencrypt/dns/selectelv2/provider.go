package selectelv2

import (
	"context"

	"github.com/go-acme/lego/v5/challenge"
	"github.com/go-acme/lego/v5/providers/dns/selectelv2"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/dynamicfields"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/i18n"
	dns2 "github.com/lucasdillmann/nginx-ignition/internal/certificate/domain/letsencrypt/dns"
)

//nolint:gosec
const (
	baseURLFieldID   = "selectelv2BaseUrl"
	usernameFieldID  = "selectelv2Username"
	passwordFieldID  = "selectelv2Password"
	projectIDFieldID = "selectelv2ProjectId"
	accountFieldID   = "selectelv2Account"
	regionFieldID    = "selectelv2Region"
)

type Provider struct{}

func (p *Provider) ID() string {
	return "SELECTELV2"
}

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateLetsencryptDnsSelectelv2Name)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns2.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          usernameFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsSelectelv2Username),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          passwordFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsSelectelv2Password),
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          projectIDFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsSelectelv2ProjectId),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          baseURLFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsSelectelv2BaseUrl),
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          accountFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsSelectelv2AccountName),
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          regionFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsSelectelv2Region),
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
	password, _ := parameters[passwordFieldID].(string)
	projectID, _ := parameters[projectIDFieldID].(string)
	baseURL, _ := parameters[baseURLFieldID].(string)
	account, _ := parameters[accountFieldID].(string)
	region, _ := parameters[regionFieldID].(string)

	cfg := selectelv2.NewDefaultConfig()
	cfg.Username = username
	cfg.Password = password
	cfg.ProjectID = projectID
	cfg.BaseURL = baseURL
	cfg.DomainName = account
	cfg.AuthRegion = region
	cfg.TTL = dns2.TTL
	cfg.PropagationTimeout = dns2.PropagationTimeout
	cfg.PollingInterval = dns2.PollingInterval

	return selectelv2.NewDNSProviderConfig(cfg)
}
