package timewebcloud

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/timewebcloud"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

const (
	authTokenFieldID = "timewebCloudAuthToken"
)

type Provider struct{}

func (p *Provider) ID() string {
	return "TIMEWEBCLOUD"
}

func (p *Provider) Name() string {
	return "Timeweb Cloud"
}

func (p *Provider) DynamicFields() []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          authTokenFieldID,
			Description: "Timeweb Cloud auth token",
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
	authToken, _ := parameters[authTokenFieldID].(string)

	cfg := &timewebcloud.Config{
		AuthToken:          authToken,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PollingInterval,
	}

	return timewebcloud.NewDNSProviderConfig(cfg)
}
