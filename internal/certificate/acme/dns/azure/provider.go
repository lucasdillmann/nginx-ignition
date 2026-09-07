package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/go-acme/lego/v5/challenge"
	"github.com/go-acme/lego/v5/providers/dns/azuredns"

	"github.com/lucasdillmann/nginx-ignition/internal/certificate/acme/dns"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/coreerror"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/dynamicfields"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/i18n"
)

//nolint:gosec
const (
	tenantFieldID       = "azureTenantId"
	subscriptionFieldID = "azureSubscriptionId"
	clientFieldID       = "azureClientId"
	clientSecretFieldID = "azureClientSecret"
	environmentFieldID  = "azureEnvironment"

	defaultRegion = "DEFAULT"
	chinaRegion   = "CHINA"
	usGovRegion   = "US_GOVERNMENT"
)

type Provider struct{}

func (p *Provider) ID() string {
	return "AZURE"
}

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateAcmeDnsAzureName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          tenantFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateAcmeDnsAzureTenantId),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          subscriptionFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateAcmeDnsAzureSubscriptionId),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          clientFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateAcmeDnsAzureClientId),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          clientSecretFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateAcmeDnsAzureClientSecret),
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:           environmentFieldID,
			Description:  i18n.M(ctx, i18n.K.CertificateAcmeDnsAzureEnvironment),
			Required:     true,
			DefaultValue: defaultRegion,
			Type:         dynamicfields.EnumType,
			EnumOptions: []dynamicfields.EnumOption{
				{
					ID: defaultRegion,
					Description: i18n.M(
						ctx,
						i18n.K.CertificateAcmeDnsAzureEnvironmentDefault,
					),
				},
				{
					ID: chinaRegion,
					Description: i18n.M(
						ctx,
						i18n.K.CertificateAcmeDnsAzureEnvironmentChina,
					),
				},
				{
					ID: usGovRegion,
					Description: i18n.M(
						ctx,
						i18n.K.CertificateAcmeDnsAzureEnvironmentUsGovernment,
					),
				},
			},
		},
	})
}

func (p *Provider) ChallengeProvider(
	ctx context.Context,
	_ []string,
	parameters map[string]any,
) (challenge.Provider, error) {
	tenantID, _ := parameters[tenantFieldID].(string)
	subscriptionID, _ := parameters[subscriptionFieldID].(string)
	clientID, _ := parameters[clientFieldID].(string)
	clientSecret, _ := parameters[clientSecretFieldID].(string)
	environment, _ := parameters[environmentFieldID].(string)

	var env cloud.Configuration
	switch environment {
	case chinaRegion:
		env = cloud.AzureChina
	case defaultRegion:
		env = cloud.AzurePublic
	case usGovRegion:
		env = cloud.AzureGovernment
	default:
		return nil, coreerror.New(
			i18n.M(ctx, i18n.K.CertificateAcmeDnsAzureErrorAzureUnknownEnvironment),
			true,
		)
	}

	cfg := azuredns.NewDefaultConfig()
	cfg.TenantID = tenantID
	cfg.SubscriptionID = subscriptionID
	cfg.ClientID = clientID
	cfg.ClientSecret = clientSecret
	cfg.Environment = env
	cfg.TTL = dns.TTL
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval

	return azuredns.NewDNSProviderConfig(cfg)
}
