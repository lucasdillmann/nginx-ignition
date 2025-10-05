package transip

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/transip"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

const (
	accountNameFieldID    = "transIpAccountName"
	privateKeyPathFieldID = "transIpPrivateKeyPath"
)

type Provider struct{}

func (p *Provider) ID() string { return "TRANSIP" }

func (p *Provider) Name() string { return "TransIP" }

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamic_fields.DynamicField{
		{
			ID:          accountNameFieldID,
			Description: "TransIP account name",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          privateKeyPathFieldID,
			Description: "TransIP private key file path",
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
	accountName, _ := parameters[accountNameFieldID].(string)
	privateKeyPath, _ := parameters[privateKeyPathFieldID].(string)

	cfg := &transip.Config{
		AccountName:        accountName,
		PrivateKeyPath:     privateKeyPath,
		TTL:                int64(dns.TTL),
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PoolingInterval,
	}

	return transip.NewDNSProviderConfig(cfg)
}
