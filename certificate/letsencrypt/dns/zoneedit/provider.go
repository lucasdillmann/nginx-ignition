package zoneedit

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/zoneedit"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

const (
	userFieldID      = "zoneeditUser"
	authTokenFieldID = "zoneeditAuthToken"
)

type Provider struct{}

func (p *Provider) ID() string {
	return "ZONEEDIT"
}

func (p *Provider) Name() string {
	return "ZoneEdit"
}

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamic_fields.DynamicField{
		{
			ID:          userFieldID,
			Description: "ZoneEdit user",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          authTokenFieldID,
			Description: "ZoneEdit auth token",
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
	user, _ := parameters[userFieldID].(string)
	authToken, _ := parameters[authTokenFieldID].(string)

	cfg := &zoneedit.Config{
		User:               user,
		AuthToken:          authToken,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PollingInterval,
	}

	return zoneedit.NewDNSProviderConfig(cfg)
}
