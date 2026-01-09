package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	route53client "github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/route53"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/coreerror"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

//nolint:gosec
const (
	region           = "us-east-1"
	accessKeyFieldID = "awsAccessKey"
	secretKeyFieldID = "awsSecretKey"
)

type Provider struct{}

func (p *Provider) ID() string {
	return "AWS_ROUTE53"
}

func (p *Provider) Name() string {
	return "AWS Route53"
}

func (p *Provider) DynamicFields() []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          accessKeyFieldID,
			Description: "AWS access key",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          secretKeyFieldID,
			Description: "AWS secret key",
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
	})
}

func (p *Provider) ChallengeProvider(
	ctx context.Context,
	domainNames []string,
	parameters map[string]any,
) (challenge.Provider, error) {
	accessKey, _ := parameters[accessKeyFieldID].(string)
	secretKey, _ := parameters[secretKeyFieldID].(string)

	hostedZoneID, err := resolveHostedZoneID(ctx, accessKey, secretKey, domainNames)
	if err != nil {
		return nil, err
	}

	cfg := route53.NewDefaultConfig()
	cfg.AccessKeyID = accessKey
	cfg.SecretAccessKey = secretKey
	cfg.HostedZoneID = *hostedZoneID
	cfg.Region = region
	cfg.MaxRetries = dns.MaxRetries
	cfg.TTL = dns.TTL
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval

	return route53.NewDNSProviderConfig(cfg)
}

func resolveHostedZoneID(
	ctx context.Context,
	accessKey, secretKey string,
	domainNames []string,
) (*string, error) {
	cfg := aws.Config{
		Credentials: credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
		Region:      region,
	}

	client := route53client.NewFromConfig(cfg)
	hostedZones, err := client.ListHostedZones(ctx, nil)
	if err != nil {
		return nil, err
	}

	for _, domainName := range domainNames {
		var bestMatch *string
		var bestMatchLength int

		for _, hostedZone := range hostedZones.HostedZones {
			hostedZoneName := strings.TrimSuffix(*hostedZone.Name, ".")

			if domainName == hostedZoneName || strings.HasSuffix(domainName, "."+hostedZoneName) {
				if len(hostedZoneName) > bestMatchLength {
					bestMatch = hostedZone.Id
					bestMatchLength = len(hostedZoneName)
				}
			}
		}

		if bestMatch != nil {
			return bestMatch, nil
		}
	}

	return nil, coreerror.New("AWS Hosted Zone ID not found from given domain names", true)
}
