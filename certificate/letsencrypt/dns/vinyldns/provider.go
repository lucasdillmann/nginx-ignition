package vinyldns

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/vinyldns"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

//nolint:gosec
const (
	accessKeyFieldID = "vinylDnsAccessKey"
	secretKeyFieldID = "vinylDnsSecretKey"
	hostFieldID      = "vinylDnsHost"
)

type Provider struct{}

func (p *Provider) ID() string {
	return "VINYLDNS"
}

func (p *Provider) Name() string {
	return "VinylDNS"
}

func (p *Provider) DynamicFields() []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          accessKeyFieldID,
			Description: "VinylDNS access key",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          secretKeyFieldID,
			Description: "VinylDNS secret key",
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          hostFieldID,
			Description: "VinylDNS host",
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
	accessKey, _ := parameters[accessKeyFieldID].(string)
	secretKey, _ := parameters[secretKeyFieldID].(string)
	host, _ := parameters[hostFieldID].(string)

	cfg := vinyldns.NewDefaultConfig()
	cfg.AccessKey = accessKey
	cfg.SecretKey = secretKey
	cfg.Host = host
	cfg.TTL = dns.TTL
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval

	return vinyldns.NewDNSProviderConfig(cfg)
}
