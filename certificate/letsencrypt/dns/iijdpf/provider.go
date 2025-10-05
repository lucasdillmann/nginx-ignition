package iijdpf

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/iijdpf"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

const (
	tokenFieldID       = "iijDpfToken"
	serviceCodeFieldID = "iijDpfServiceCode"
	endpointFieldID    = "iijDpfEndpoint"
)

type Provider struct{}

func (p *Provider) ID() string { return "IIJ_DPF" }

func (p *Provider) Name() string { return "IIJ DNS Platform Service" }

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamic_fields.DynamicField{
		{
			ID:          tokenFieldID,
			Description: "IIJ DPF API token",
			Required:    true,
			Sensitive:   true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          serviceCodeFieldID,
			Description: "IIJ DPF DPM service code",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          endpointFieldID,
			Description: "IIJ DPF API endpoint",
			Type:        dynamic_fields.SingleLineTextType,
		},
	})
}

func (p *Provider) ChallengeProvider(
	_ context.Context,
	_ []string,
	parameters map[string]any,
) (challenge.Provider, error) {
	token, _ := parameters[tokenFieldID].(string)
	serviceCode, _ := parameters[serviceCodeFieldID].(string)
	endpoint, _ := parameters[endpointFieldID].(string)

	cfg := &iijdpf.Config{
		Token:              token,
		ServiceCode:        serviceCode,
		Endpoint:           endpoint,
		TTL:                dns.TTL,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PoolingInterval,
	}

	return iijdpf.NewDNSProviderConfig(cfg)
}
