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
	tenantID       = "azureTenantId"
	subscriptionID = "azureSubscriptionId"
	clientID       = "azureClientId"
	clientSecretID = "azureClientSecret"
	environmentID  = "azureEnvironment"
)

type Provider struct{}

func (p *Provider) ID() string {
	return "AZURE"
}

func (p *Provider) Name() string {
	return "Azure DNS"
}

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamic_fields.DynamicField{
		{
			ID:          tenantID,
			Description: "Azure tenant ID (for the DNS challenge)",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          subscriptionID,
			Description: "Azure subscription ID",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          clientID,
			Description: "Azure client ID",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          clientSecretID,
			Description: "Azure client secret",
			Required:    true,
			Sensitive:   true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          environmentID,
			Description: "Azure environment",
			Required:    true,
			Type:        dynamic_fields.EnumType,
			EnumOptions: &[]*dynamic_fields.EnumOption{
				{ID: "DEFAULT", Description: "Azure (default)"},
				{ID: "CHINA", Description: "China"},
				{ID: "US_GOVERNMENT", Description: "US Government"},
			},
		},
	})
}

func (p *Provider) ChallengeProvider(
	_ context.Context,
	_ []string,
	parameters map[string]any,
) (challenge.Provider, error) {
	tenantId, _ := parameters[tenantID].(string)
	subscriptionId, _ := parameters[subscriptionID].(string)
	clientId, _ := parameters[clientID].(string)
	clientSecret, _ := parameters[clientSecretID].(string)
	environment, _ := parameters[environmentID].(string)

	var env cloud.Configuration
	switch environment {
	case "CHINA":
		env = cloud.AzureChina
	case "DEFAULT":
		env = cloud.AzurePublic
	case "US_GOVERNMENT":
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
