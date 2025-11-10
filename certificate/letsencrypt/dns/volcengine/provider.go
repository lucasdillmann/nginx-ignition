package volcengine

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/volcengine"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

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

func (p *Provider) DynamicFields() []*dynamicfields.DynamicField {
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

	cfg := &volcengine.Config{
		AccessKey:          accessKey,
		SecretKey:          secretKey,
		Region:             region,
		TTL:                dns.TTL,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PollingInterval,
	}

	return volcengine.NewDNSProviderConfig(cfg)
}
