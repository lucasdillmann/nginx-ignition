package httpreq

import (
	"context"
	"net/url"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/httpreq"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/core_error"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

const (
	endpointFieldID = "httpReqEndpoint"
	usernameFieldID = "httpReqUsername"
	passwordFieldID = "httpReqPassword"
	modeFieldID     = "httpReqMode"
)

type Provider struct{}

func (p *Provider) ID() string { return "HTTP_REQUEST" }

func (p *Provider) Name() string { return "HTTP request" }

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamic_fields.DynamicField{
		{
			ID:          endpointFieldID,
			Description: "HTTP endpoint URL",
			Required:    true,
			Type:        dynamic_fields.URLType,
		},
		{
			ID:          usernameFieldID,
			Description: "HTTP basic auth username",
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          passwordFieldID,
			Description: "HTTP basic auth password",
			Sensitive:   true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          modeFieldID,
			Description: "Raw mode",
			Type:        dynamic_fields.BooleanType,
		},
	})
}

func (p *Provider) ChallengeProvider(_ context.Context, _ []string, parameters map[string]any) (challenge.Provider, error) {
	endpointStr, _ := parameters[endpointFieldID].(string)
	username, _ := parameters[usernameFieldID].(string)
	password, _ := parameters[passwordFieldID].(string)
	rawMode, _ := parameters[modeFieldID].(bool)

	mode := ""
	if rawMode {
		mode = "RAW"
	}

	endpoint, err := url.Parse(endpointStr)
	if err != nil {
		return nil, core_error.New("Invalid endpoint URL", true)
	}

	cfg := &httpreq.Config{
		Endpoint:           endpoint,
		Mode:               mode,
		Username:           username,
		Password:           password,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PoolingInterval,
	}

	return httpreq.NewDNSProviderConfig(cfg)
}
