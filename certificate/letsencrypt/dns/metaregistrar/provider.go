package metaregistrar

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/metaregistrar"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

const (
	apiTokenFieldID = "metaregistrarApiToken"
)

type Provider struct{}

func (p *Provider) ID() string { return "METAREGISTRAR" }

func (p *Provider) Name() string { return "Metaregistrar" }

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamic_fields.DynamicField{
		{
			ID:          apiTokenFieldID,
			Description: "Metaregistrar API token",
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
	apiToken, _ := parameters[apiTokenFieldID].(string)

	cfg := &metaregistrar.Config{
		APIToken:           apiToken,
		TTL:                dns.TTL,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PollingInterval,
	}

	return metaregistrar.NewDNSProviderConfig(cfg)
}
