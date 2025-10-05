package edgedns

import (
	"context"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/edgegrid"
	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/edgedns"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

const (
	hostFieldID         = "edgeDnsHost"
	clientTokenFieldID  = "edgeDnsClientToken"
	clientSecretFieldID = "edgeDnsClientSecret"
	accessTokenFieldID  = "edgeDnsAccessToken"
	accountKeyFieldID   = "edgeDnsAccountKey"
)

type Provider struct{}

func (p *Provider) ID() string { return "EDGEDNS" }

func (p *Provider) Name() string { return "Akamai EdgeDNS" }

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamic_fields.DynamicField{
		{
			ID:          hostFieldID,
			Description: "Akamai EdgeGrid host",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          clientTokenFieldID,
			Description: "Akamai EdgeGrid client token",
			Required:    true,
			Sensitive:   true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          clientSecretFieldID,
			Description: "Akamai EdgeGrid client secret",
			Required:    true,
			Sensitive:   true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          accessTokenFieldID,
			Description: "Akamai EdgeGrid access token",
			Required:    true,
			Sensitive:   true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          accountKeyFieldID,
			Description: "Akamai EdgeGrid account key",
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
	host, _ := parameters[hostFieldID].(string)
	clientToken, _ := parameters[clientTokenFieldID].(string)
	clientSecret, _ := parameters[clientSecretFieldID].(string)
	accessToken, _ := parameters[accessTokenFieldID].(string)
	accountKey, _ := parameters[accountKeyFieldID].(string)

	cfg := &edgedns.Config{
		Config: &edgegrid.Config{
			Host:         host,
			ClientToken:  clientToken,
			ClientSecret: clientSecret,
			AccessToken:  accessToken,
			AccountKey:   accountKey,
		},
		TTL:                dns.TTL,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PoolingInterval,
	}

	return edgedns.NewDNSProviderConfig(cfg)
}
