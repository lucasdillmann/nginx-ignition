package directadmin

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/directadmin"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

const (
	hostFieldID     = "directAdminHost"
	passwordFieldID = "directAdminPassword"
	userFieldID     = "directAdminUsername"
	zoneNameFieldID = "directAdminZoneName"
)

type Provider struct{}

func (p *Provider) ID() string { return "DIRECTADMIN" }

func (p *Provider) Name() string { return "DirectAdmin" }

func (p *Provider) DynamicFields() []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          hostFieldID,
			Description: "DirectAdmin base URL",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          userFieldID,
			Description: "DirectAdmin username",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          passwordFieldID,
			Description: "DirectAdmin password",
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          zoneNameFieldID,
			Description: "DirectAdmin zone name",
			Type:        dynamicfields.SingleLineTextType,
		},
	})
}

func (p *Provider) ChallengeProvider(
	_ context.Context,
	_ []string,
	parameters map[string]any,
) (challenge.Provider, error) {
	host, _ := parameters[hostFieldID].(string)
	user, _ := parameters[userFieldID].(string)
	password, _ := parameters[passwordFieldID].(string)
	zoneName, _ := parameters[zoneNameFieldID].(string)

	cfg := directadmin.NewDefaultConfig()
	cfg.BaseURL = host
	cfg.Username = user
	cfg.Password = password
	cfg.ZoneName = zoneName
	cfg.TTL = dns.TTL
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval
	cfg.HTTPClient = nil

	return directadmin.NewDNSProviderConfig(cfg)
}
