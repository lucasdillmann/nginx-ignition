package volcengine

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/volcengine"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

//nolint:gosec
const (
	accessKeyFieldID = "volcengineAccessKey"
	secretKeyFieldID = "volcengineSecretKey"
	regionFieldID    = "volcengineRegion"
)

type Provider struct{}

func (p *Provider) ID() string {
	return "VOLCENGINE"
}

func (p *Provider) Name() string {
	return "Volcengine"
}

func (p *Provider) DynamicFields() []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          accessKeyFieldID,
			Description: "Volcengine access key",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          secretKeyFieldID,
			Description: "Volcengine secret key",
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          regionFieldID,
			Description: "Volcengine region",
			Type:        dynamicfields.SingleLineTextType,
		},
	})
}

func (p *Provider) ChallengeProvider(
	_ context.Context,
	_ []string,
	parameters map[string]any,
) (challenge.Provider, error) {
	accessKey, _ := parameters[accessKeyFieldID].(string)
	secretKey, _ := parameters[secretKeyFieldID].(string)
	region, _ := parameters[regionFieldID].(string)

	cfg := volcengine.NewDefaultConfig()
	cfg.AccessKey = accessKey
	cfg.SecretKey = secretKey
	cfg.Region = region
	cfg.TTL = dns.TTL
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval

	return volcengine.NewDNSProviderConfig(cfg)
}
