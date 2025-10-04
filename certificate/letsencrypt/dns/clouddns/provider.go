package clouddns

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/clouddns"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

const (
	clientIDFieldID = "cloudDnsClientId"
	emailFieldID    = "cloudDnsEmail"
	passwordFieldID = "cloudDnsPassword"
)

type Provider struct{}

func (p *Provider) ID() string { return "CLOUDDNS" }

func (p *Provider) Name() string { return "CloudDNS" }

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamic_fields.DynamicField{
		{
			ID:          clientIDFieldID,
			Description: "CloudDNS client ID",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          emailFieldID,
			Description: "CloudDNS email",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          passwordFieldID,
			Description: "CloudDNS password",
			Required:    true,
			Sensitive:   true,
			Type:        dynamic_fields.SingleLineTextType,
		},
	})
}

func (p *Provider) ChallengeProvider(_ context.Context, _ []string, parameters map[string]any) (challenge.Provider, error) {
	clientID, _ := parameters[clientIDFieldID].(string)
	email, _ := parameters[emailFieldID].(string)
	password, _ := parameters[passwordFieldID].(string)

	cfg := &clouddns.Config{
		ClientID:           clientID,
		Email:              email,
		Password:           password,
		TTL:                dns.TTL,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PoolingInterval,
	}

	return clouddns.NewDNSProviderConfig(cfg)
}
