package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/azuredns"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/coreerror"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
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

func (p *Provider) Name() string {
	return "Azure"
}

func (p *Provider) DynamicFields() []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          tenantFieldID,
			Description: "Azure tenant ID",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          subscriptionFieldID,
			Description: "Azure subscription ID",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          clientFieldID,
			Description: "Azure client ID",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          clientSecretFieldID,
			Description: "Azure client secret",
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:           environmentFieldID,
			Description:  "Azure environment",
			Required:     true,
			DefaultValue: defaultRegion,
			Type:         dynamicfields.EnumType,
			EnumOptions: []dynamicfields.EnumOption{
				{
					ID:          defaultRegion,
					Description: "Azure (default)",
				},
				{
					ID:          chinaRegion,
					Description: "China",
				},
				{
					ID:          usGovRegion,
					Description: "US Government",
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
		return nil, coreerror.New(i18n.M(ctx, i18n.K.CertificateErrorAzureUnknownEnvironment), true)
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
