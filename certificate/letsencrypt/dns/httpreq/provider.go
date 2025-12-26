package httpreq

import (
	"context"
	"net/url"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/httpreq"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/coreerror"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
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

func (p *Provider) DynamicFields() []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          endpointFieldID,
			Description: "HTTP endpoint URL",
			Required:    true,
			Type:        dynamicfields.URLType,
		},
		{
			ID:          usernameFieldID,
			Description: "HTTP basic auth username",
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          passwordFieldID,
			Description: "HTTP basic auth password",
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          modeFieldID,
			Description: "Raw mode",
			Type:        dynamicfields.BooleanType,
		},
	})
}

func (p *Provider) ChallengeProvider(
	_ context.Context,
	_ []string,
	parameters map[string]any,
) (challenge.Provider, error) {
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
		return nil, coreerror.New("Invalid endpoint URL", true)
	}

	cfg := &httpreq.Config{
		Endpoint:           endpoint,
		Mode:               mode,
		Username:           username,
		Password:           password,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PollingInterval,
	}

	return httpreq.NewDNSProviderConfig(cfg)
}
