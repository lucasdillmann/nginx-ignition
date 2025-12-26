package yandex360

import (
	"context"
	"strconv"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/yandex360"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

//nolint:gosec
const (
	oauthTokenFieldID = "yandex360OAuthToken"
	orgIDFieldID      = "yandex360OrgId"
)

type Provider struct{}

func (p *Provider) ID() string {
	return "YANDEX360"
}

func (p *Provider) Name() string {
	return "Yandex 360"
}

func (p *Provider) DynamicFields() []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          oauthTokenFieldID,
			Description: "Yandex 360 OAuth token",
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          orgIDFieldID,
			Description: "Yandex 360 organization ID",
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
	oauthToken, _ := parameters[oauthTokenFieldID].(string)
	orgIDStr, _ := parameters[orgIDFieldID].(string)

	orgID, err := strconv.ParseInt(orgIDStr, 10, 64)
	if err != nil {
		return nil, err
	}

	cfg := &yandex360.Config{
		OAuthToken:         oauthToken,
		OrgID:              orgID,
		TTL:                dns.TTL,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PollingInterval,
	}

	return yandex360.NewDNSProviderConfig(cfg)
}
