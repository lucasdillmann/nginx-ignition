package bluecat

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/bluecat"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

const (
	baseURLFieldID    = "blueCatBaseURL"
	usernameFieldID   = "blueCatUsername"
	passwordFieldID   = "blueCatPassword"
	configNameFieldID = "blueCatConfigName"
	dnsViewFieldID    = "blueCatDNSView"
)

type Provider struct{}

func (p *Provider) ID() string { return "BLUECAT" }

func (p *Provider) Name() string { return "BlueCat" }

func (p *Provider) DynamicFields() []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          baseURLFieldID,
			Description: "BlueCat base URL",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          usernameFieldID,
			Description: "BlueCat username",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          passwordFieldID,
			Description: "BlueCat password",
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          configNameFieldID,
			Description: "BlueCat configuration name",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          dnsViewFieldID,
			Description: "BlueCat DNS view",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
	})
}

func (p *Provider) ChallengeProvider(
	_ context.Context,
	_ []string,
	parameters map[string]any,
) (challenge.Provider, error) {
	baseURL, _ := parameters[baseURLFieldID].(string)
	username, _ := parameters[usernameFieldID].(string)
	password, _ := parameters[passwordFieldID].(string)
	configName, _ := parameters[configNameFieldID].(string)
	dnsView, _ := parameters[dnsViewFieldID].(string)

	cfg := &bluecat.Config{
		BaseURL:            baseURL,
		UserName:           username,
		Password:           password,
		ConfigName:         configName,
		DNSView:            dnsView,
		TTL:                dns.TTL,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PollingInterval,
	}

	return bluecat.NewDNSProviderConfig(cfg)
}
