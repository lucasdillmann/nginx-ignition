package mydnsjp

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/mydnsjp"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

const (
	masterIDFieldID = "myDnsJpMasterId"
	passwordFieldID = "myDnsJpPassword"
)

type Provider struct{}

func (p *Provider) ID() string { return "MYDNS_JP" }

func (p *Provider) Name() string { return "MyDNS.jp" }

func (p *Provider) DynamicFields() []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          masterIDFieldID,
			Description: "MyDNS.jp master ID",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          passwordFieldID,
			Description: "MyDNS.jp password",
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
	})
}

func (p *Provider) ChallengeProvider(
	_ context.Context,
	_ []string,
	parameters map[string]any,
) (challenge.Provider, error) {
	masterID, _ := parameters[masterIDFieldID].(string)
	password, _ := parameters[passwordFieldID].(string)

	cfg := mydnsjp.NewDefaultConfig()
	cfg.MasterID = masterID
	cfg.Password = password
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval

	return mydnsjp.NewDNSProviderConfig(cfg)
}
