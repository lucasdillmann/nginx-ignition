package designate

import (
	"context"
	"os"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/designate"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

const (
	authURLFieldID   = "designateAuthUrl"
	regionFieldID    = "designateRegionName"
	tenantFieldID    = "designateTenantName"
	usernameFieldID  = "designateUsername"
	passwordFieldID  = "designatePassword"
	projectIDFieldID = "designateProjectId"
)

type Provider struct{}

func (p *Provider) ID() string { return "OPENSTACK_DESIGNATE" }

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateLetsencryptDnsDesignateName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          authURLFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsDesignateAuthUrl),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          regionFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsDesignateRegionName),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          tenantFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsDesignateTenantName),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          usernameFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsDesignateUsername),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          passwordFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsDesignatePassword),
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          projectIDFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsDesignateProjectId),
			Type:        dynamicfields.SingleLineTextType,
		},
	})
}

func (p *Provider) ChallengeProvider(
	_ context.Context,
	_ []string,
	parameters map[string]any,
) (challenge.Provider, error) {
	authURL, _ := parameters[authURLFieldID].(string)
	region, _ := parameters[regionFieldID].(string)
	tenant, _ := parameters[tenantFieldID].(string)
	username, _ := parameters[usernameFieldID].(string)
	password, _ := parameters[passwordFieldID].(string)
	projectID, _ := parameters[projectIDFieldID].(string)

	envVars := map[string]string{
		"OS_AUTH_URL":    authURL,
		"OS_USERNAME":    username,
		"OS_PASSWORD":    password,
		"OS_TENANT_NAME": tenant,
		"OS_REGION_NAME": region,
		"OS_PROJECT_ID":  projectID,
	}

	previousValues := make(map[string]string)
	for key, value := range envVars {
		previousValues[key] = os.Getenv(key)
		_ = os.Setenv(key, value)
	}

	defer func() {
		for key, value := range previousValues {
			if value == "" {
				_ = os.Unsetenv(key)
			} else {
				_ = os.Setenv(key, value)
			}
		}
	}()

	cfg := designate.NewDefaultConfig()
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval
	cfg.TTL = dns.TTL

	return designate.NewDNSProviderConfig(cfg)
}
