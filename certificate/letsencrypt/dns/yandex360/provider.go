package yandex360

import (
	"context"
	"strconv"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/yandex360"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

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

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamic_fields.DynamicField{
		{
			ID:          oauthTokenFieldID,
			Description: "Yandex 360 OAuth token",
			Required:    true,
			Sensitive:   true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          orgIDFieldID,
			Description: "Yandex 360 organization ID",
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
		PollingInterval:    dns.PoolingInterval,
	}

	return yandex360.NewDNSProviderConfig(cfg)
}
