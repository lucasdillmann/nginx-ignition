package edgeone

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/edgeone"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

//nolint:gosec
const (
	secretIDFieldID  = "edgeOneSecretID"
	secretKeyFieldID = "edgeOneSecretKey"
	regionFieldID    = "edgeOneRegion"
)

type Provider struct{}

func (p *Provider) ID() string { return "EDGEONE" }

func (p *Provider) Name() string { return "Tencent EdgeOne" }

func (p *Provider) DynamicFields() []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          secretIDFieldID,
			Description: "Tencent Cloud secret ID",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          secretKeyFieldID,
			Description: "Tencent Cloud secret key",
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          regionFieldID,
			Description: "Tencent Cloud region",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
	})
}

func (p *Provider) ChallengeProvider(
	_ context.Context,
	_ []string,
	parameters map[string]any,
) (challenge.Provider, error) {
	secretID, _ := parameters[secretIDFieldID].(string)
	secretKey, _ := parameters[secretKeyFieldID].(string)
	region, _ := parameters[regionFieldID].(string)

	cfg := edgeone.NewDefaultConfig()
	cfg.SecretID = secretID
	cfg.SecretKey = secretKey
	cfg.Region = region
	cfg.TTL = dns.TTL
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval

	return edgeone.NewDNSProviderConfig(cfg)
}
