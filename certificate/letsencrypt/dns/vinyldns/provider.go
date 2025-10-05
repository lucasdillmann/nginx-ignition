package vinyldns

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/vinyldns"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

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

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamic_fields.DynamicField{
		{
			ID:          accessKeyFieldID,
			Description: "VinylDNS access key",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          secretKeyFieldID,
			Description: "VinylDNS secret key",
			Required:    true,
			Sensitive:   true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          hostFieldID,
			Description: "VinylDNS host",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
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

	cfg := &vinyldns.Config{
		AccessKey:          accessKey,
		SecretKey:          secretKey,
		Host:               host,
		TTL:                dns.TTL,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PollingInterval,
	}

	return vinyldns.NewDNSProviderConfig(cfg)
}
