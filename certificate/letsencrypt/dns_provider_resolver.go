package letsencrypt

import (
	"context"
	"dillmann.com.br/nginx-ignition/core/common/core_error"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	route53client "github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/azuredns"
	"github.com/go-acme/lego/v4/providers/dns/cloudflare"
	"github.com/go-acme/lego/v4/providers/dns/gcloud"
	"github.com/go-acme/lego/v4/providers/dns/route53"
	"strings"
	"time"
)

const (
	ttl                = 60
	maxRetries         = 3
	propagationTimeout = 180 * time.Second
	poolingInterval    = 1 * time.Second
)

func resolveDnsProvider(domainNames []string, parameters map[string]any) (challenge.Provider, error) {
	providerId, _ := parameters[dnsProvider.ID].(string)

	switch providerId {
	case awsRoute53Id:
		return buildAwsRoute53Provider(domainNames, parameters)
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

func buildAwsRoute53Provider(domainNames []string, parameters map[string]any) (challenge.Provider, error) {
	accessKey, _ := parameters[awsAccessKey.ID].(string)
	secretKey, _ := parameters[awsSecretKey.ID].(string)

	hostedZoneId, err := resolveAwsRoute53HostedZoneID(accessKey, secretKey, domainNames)
	if err != nil {
		return nil, err
	}

	cfg := &route53.Config{
		AccessKeyID:        accessKey,
		SecretAccessKey:    secretKey,
		HostedZoneID:       *hostedZoneId,
		MaxRetries:         maxRetries,
		TTL:                ttl,
		PropagationTimeout: propagationTimeout,
		PollingInterval:    poolingInterval,
		Region:             "us-east-1",
	}

	return route53.NewDNSProviderConfig(cfg)
}

func resolveAwsRoute53HostedZoneID(accessKey, secretKey string, domainNames []string) (*string, error) {
	cfg := aws.Config{
		Credentials: credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
		Region:      "us-east-1",
	}

	client := route53client.NewFromConfig(cfg)
	hostedZones, err := client.ListHostedZones(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	for _, domainName := range domainNames {
		for _, hostedZone := range hostedZones.HostedZones {
			hostedZoneName := strings.TrimSuffix(*hostedZone.Name, ".")
			if strings.HasSuffix(domainName, hostedZoneName) {
				return hostedZone.Id, nil
			}
		}
	}

	return nil, core_error.New("AWS Hosted Zone ID not found from given domain names", true)
}

func buildCloudflareProvider(parameters map[string]any) (challenge.Provider, error) {
	apiToken, _ := parameters[cloudflareApiToken.ID].(string)

	cfg := &cloudflare.Config{
		AuthToken: apiToken,
	}

	return cloudflare.NewDNSProviderConfig(cfg)
}

func buildGoogleCloudProvider(parameters map[string]any) (challenge.Provider, error) {
	privateKey, _ := parameters[googleCloudPrivateKey.ID].(string)
	privateKeyBytes := []byte(privateKey)
	return gcloud.NewDNSProviderServiceAccountKey(privateKeyBytes)
}

func buildAzureProvider(parameters map[string]any) (challenge.Provider, error) {
	tenantId, _ := parameters[azureTenantId.ID].(string)
	subscriptionId, _ := parameters[azureSubscriptionId.ID].(string)
	clientId, _ := parameters[azureClientId.ID].(string)
	clientSecret, _ := parameters[azureClientSecret.ID].(string)
	environmentID, _ := parameters[azureEnvironment.ID].(string)

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
