package cpanel

import (
	"context"

	"github.com/aws/smithy-go/ptr"
	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/cpanel"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

const (
	hostFieldID  = "cpanelHost"
	tokenFieldID = "cpanelToken"
	userFieldID  = "cpanelUsername"
	modeFieldID  = "cpanelMode"
)

type Provider struct{}

func (p *Provider) ID() string { return "CPANEL" }

func (p *Provider) Name() string { return "cPanel" }

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamic_fields.DynamicField{
		{
			ID:          hostFieldID,
			Description: "cPanel base URL",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          userFieldID,
			Description: "cPanel username",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          tokenFieldID,
			Description: "cPanel API token",
			Required:    true,
			Sensitive:   true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          modeFieldID,
			Description: "cPanel mode",
			HelpText:    ptr.String("Defaults to cpanel when left empty"),
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
	user, _ := parameters[userFieldID].(string)
	token, _ := parameters[tokenFieldID].(string)
	mode, _ := parameters[modeFieldID].(string)

	if mode == "" {
		mode = "cpanel"
	}

	cfg := &cpanel.Config{
		BaseURL:            host,
		Username:           user,
		Token:              token,
		Mode:               mode,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PollingInterval,
	}

	return cpanel.NewDNSProviderConfig(cfg)
}
