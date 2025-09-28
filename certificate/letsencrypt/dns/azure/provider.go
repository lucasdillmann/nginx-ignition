package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/azuredns"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/core_error"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

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

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamic_fields.DynamicField{
		{
			ID:          tenantFieldID,
			Description: "Azure tenant ID",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          subscriptionFieldID,
			Description: "Azure subscription ID",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          clientFieldID,
			Description: "Azure client ID",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          clientSecretFieldID,
			Description: "Azure client secret",
			Required:    true,
			Sensitive:   true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          environmentFieldID,
			Description: "Azure environment",
			Required:    true,
			Type:        dynamic_fields.EnumType,
			EnumOptions: &[]*dynamic_fields.EnumOption{
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
	_ context.Context,
	_ []string,
	parameters map[string]any,
) (challenge.Provider, error) {
	tenantId, _ := parameters[tenantFieldID].(string)
	subscriptionId, _ := parameters[subscriptionFieldID].(string)
	clientId, _ := parameters[clientFieldID].(string)
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
		return nil, core_error.New("Unknown Azure environment", true)
	}

	cfg := &azuredns.Config{
		TenantID:           tenantId,
		SubscriptionID:     subscriptionId,
		ClientID:           clientId,
		ClientSecret:       clientSecret,
		Environment:        env,
		TTL:                dns.TTL,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PoolingInterval,
	}

	return azuredns.NewDNSProviderConfig(cfg)
}
