package ibmcloud

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/ibmcloud"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

const (
	usernameFieldID = "ibmCloudUsername"
	apiKeyFieldID   = "ibmCloudApiKey"
)

type Provider struct{}

func (p *Provider) ID() string { return "IBM_CLOUD" }

func (p *Provider) Name() string { return "IBM Cloud (SoftLayer)" }

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamic_fields.DynamicField{
		{
			ID:          usernameFieldID,
			Description: "IBM Cloud (SoftLayer) username",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          apiKeyFieldID,
			Description: "IBM Cloud (SoftLayer) API key",
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
	user, _ := parameters[usernameFieldID].(string)
	key, _ := parameters[apiKeyFieldID].(string)

	cfg := &ibmcloud.Config{
		Username:           user,
		APIKey:             key,
		TTL:                dns.TTL,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PollingInterval,
	}

	return ibmcloud.NewDNSProviderConfig(cfg)
}
