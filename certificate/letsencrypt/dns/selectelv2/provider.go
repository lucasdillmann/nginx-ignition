package selectelv2

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/selectelv2"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
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
	return i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsSelectelv2Name)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          usernameFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsSelectelv2Username),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          passwordFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsSelectelv2Password),
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          projectIDFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsSelectelv2ProjectId),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          baseURLFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsSelectelv2BaseUrl),
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          accountFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsSelectelv2AccountName),
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          regionFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsSelectelv2Region),
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
	cfg.TTL = dns.TTL
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval

	return selectelv2.NewDNSProviderConfig(cfg)
}
