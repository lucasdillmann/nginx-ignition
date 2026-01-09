package axelname

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/axelname"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

const (
	nicknameFieldID = "axelNameNickname"
	tokenFieldID    = "axelNameToken"
)

type Provider struct{}

func (p *Provider) ID() string { return "AXEL_NAME" }

func (p *Provider) Name() string { return "Axel Name" }

func (p *Provider) DynamicFields() []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          nicknameFieldID,
			Description: "Axel name nickname",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          tokenFieldID,
			Description: "Axel name token",
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
	nickname, _ := parameters[nicknameFieldID].(string)
	token, _ := parameters[tokenFieldID].(string)

	cfg := axelname.NewDefaultConfig()
	cfg.Nickname = nickname
	cfg.Token = token
	cfg.TTL = dns.TTL
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval

	return axelname.NewDNSProviderConfig(cfg)
}
