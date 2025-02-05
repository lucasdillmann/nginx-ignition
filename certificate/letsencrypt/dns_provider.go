package letsencrypt

import (
	"dillmann.com.br/nginx-ignition/core/common/core_error"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/azuredns"
	"github.com/go-acme/lego/v4/providers/dns/cloudflare"
	"github.com/go-acme/lego/v4/providers/dns/gcloud"
	"github.com/go-acme/lego/v4/providers/dns/route53"
	"time"
)

const (
	ttl                = 60
	maxRetries         = 3
	propagationTimeout = 180 * time.Second
	poolingInterval    = 2 * time.Second
)

func resolveDnsProvider(parameters map[string]interface{}) (challenge.Provider, error) {
	providerId, casted := parameters[dnsProvider.ID].(string)
	if !casted {
		return nil, core_error.New("DNS provider ID is missing", true)
	}

	switch providerId {
	case awsRoute53Id:
		return buildAwsRoute53Provider(parameters)
	case cloudflareId:
		return buildCloudflareProvider(parameters)
	case googleCloudId:
		return buildGoogleCloudProvider(parameters)
	case azureId:
		return buildAzureProvider(parameters)
	default:
		return nil, core_error.New("Unknown DNS provider", true)
	}
}

func buildAwsRoute53Provider(parameters map[string]interface{}) (challenge.Provider, error) {
	accessKey, casted := parameters[awsAccessKey.ID].(string)
	if !casted {
		return nil, core_error.New("AWS access key is missing", true)
	}

	secretKey, casted := parameters[awsSecretKey.ID].(string)
	if !casted {
		return nil, core_error.New("AWS secret key is missing", true)
	}

	hostedZoneID, casted := parameters[awsHostedZoneID.ID].(string)
	if !casted {
		return nil, core_error.New("AWS hosted zone ID is missing", true)
	}

	cfg := &route53.Config{
		AccessKeyID:              accessKey,
		SecretAccessKey:          secretKey,
		HostedZoneID:             hostedZoneID,
		MaxRetries:               maxRetries,
		WaitForRecordSetsChanged: true,
		TTL:                      ttl,
		PropagationTimeout:       propagationTimeout,
		PollingInterval:          poolingInterval,
	}

	return route53.NewDNSProviderConfig(cfg)
}

func buildCloudflareProvider(parameters map[string]interface{}) (challenge.Provider, error) {
	apiToken, casted := parameters[cloudflareApiToken.ID].(string)
	if !casted {
		return nil, core_error.New("Cloudflare API token is missing", true)
	}

	cfg := &cloudflare.Config{
		AuthToken: apiToken,
	}

	return cloudflare.NewDNSProviderConfig(cfg)
}

func buildGoogleCloudProvider(parameters map[string]interface{}) (challenge.Provider, error) {
	privateKey, casted := parameters[googleCloudPrivateKey.ID].(string)
	if !casted {
		return nil, core_error.New("Google Cloud private key is missing", true)
	}

	privateKeyBytes := []byte(privateKey)
	return gcloud.NewDNSProviderServiceAccountKey(privateKeyBytes)
}

func buildAzureProvider(parameters map[string]interface{}) (challenge.Provider, error) {
	tenantId, casted := parameters[azureTenantId.ID].(string)
	if !casted {
		return nil, core_error.New("Azure tenant ID is missing", true)
	}

	subscriptionId, casted := parameters[azureSubscriptionId.ID].(string)
	if !casted {
		return nil, core_error.New("Azure subscription ID is missing", true)
	}

	clientId, casted := parameters[azureClientId.ID].(string)
	if !casted {
		return nil, core_error.New("Azure client ID is missing", true)
	}

	clientSecret, casted := parameters[azureClientSecret.ID].(string)
	if !casted {
		return nil, core_error.New("Azure client secret is missing", true)
	}

	environmentID, casted := parameters[azureEnvironment.ID].(string)
	if !casted {
		return nil, core_error.New("Azure environment is missing", true)
	}

	var environment cloud.Configuration
	switch environmentID {
	case "CHINA":
		environment = cloud.AzureChina
	case "DEFAULT":
		environment = cloud.AzurePublic
	case "US_GOVERNMENT":
		environment = cloud.AzureGovernment
	default:
		return nil, core_error.New("Unknown Azure environment", true)
	}

	cfg := &azuredns.Config{
		TenantID:           tenantId,
		SubscriptionID:     subscriptionId,
		ClientID:           clientId,
		ClientSecret:       clientSecret,
		Environment:        environment,
		TTL:                ttl,
		PropagationTimeout: propagationTimeout,
		PollingInterval:    poolingInterval,
	}

	return azuredns.NewDNSProviderConfig(cfg)
}
