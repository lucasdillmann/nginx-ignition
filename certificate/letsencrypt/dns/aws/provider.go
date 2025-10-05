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
	"dillmann.com.br/nginx-ignition/core/common/core_error"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

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

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamic_fields.DynamicField{
		{
			ID:          accessKeyFieldID,
			Description: "AWS access key",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          secretKeyFieldID,
			Description: "AWS secret key",
			Required:    true,
			Sensitive:   true,
			Type:        dynamic_fields.SingleLineTextType,
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

	hostedZoneId, err := resolveHostedZoneID(ctx, accessKey, secretKey, domainNames)
	if err != nil {
		return nil, err
	}

	cfg := &route53.Config{
		AccessKeyID:        accessKey,
		SecretAccessKey:    secretKey,
		HostedZoneID:       *hostedZoneId,
		Region:             region,
		MaxRetries:         dns.MaxRetries,
		TTL:                dns.TTL,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PollingInterval,
	}

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

	return nil, core_error.New("AWS Hosted Zone ID not found from given domain names", true)
}
