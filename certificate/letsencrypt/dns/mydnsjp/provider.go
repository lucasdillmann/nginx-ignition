package mydnsjp

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/mydnsjp"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

const (
	masterIDFieldID = "myDnsJpMasterId"
	passwordFieldID = "myDnsJpPassword"
)

type Provider struct{}

func (p *Provider) ID() string { return "MYDNS_JP" }

func (p *Provider) Name() string { return "MyDNS.jp" }

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamic_fields.DynamicField{
		{
			ID:          masterIDFieldID,
			Description: "MyDNS.jp master ID",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          passwordFieldID,
			Description: "MyDNS.jp password",
			Required:    true,
			Sensitive:   true,
			Type:        dynamic_fields.SingleLineTextType,
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

	cfg := &mydnsjp.Config{
		MasterID:           masterID,
		Password:           password,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PollingInterval,
	}

	return mydnsjp.NewDNSProviderConfig(cfg)
}
