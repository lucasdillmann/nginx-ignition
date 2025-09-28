package pdns

import (
	"context"
	"errors"
	"net/url"
	"strconv"

	"github.com/aws/smithy-go/ptr"
	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/pdns"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

const (
	apiKeyFieldID     = "powerDnsApiKey"
	hostURLFieldID    = "powerDnsHostUtl"
	serverNameFieldID = "powerDnsServerName"
	apiVersionFieldID = "powerDnsApiVersion"
)

type Provider struct{}

func (p *Provider) ID() string { return "POWERDNS" }

func (p *Provider) Name() string { return "PowerDNS" }

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamic_fields.DynamicField{
		{
			ID:          apiKeyFieldID,
			Description: "PowerDNS API key",
			Required:    true,
			Sensitive:   true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          hostURLFieldID,
			Description: "PowerDNS host URL",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          serverNameFieldID,
			Description: "PowerDNS server name",
			HelpText:    ptr.String("Defaults to localhost when left empty"),
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          apiVersionFieldID,
			Description: "PowerDNS API version",
			HelpText:    ptr.String("Defaults to auto-detection when left empty"),
			Type:        dynamic_fields.SingleLineTextType,
		},
	})
}

func (p *Provider) ChallengeProvider(_ context.Context, _ []string, parameters map[string]any) (challenge.Provider, error) {
	apiKey, _ := parameters[apiKeyFieldID].(string)
	hostURLStr, _ := parameters[hostURLFieldID].(string)
	serverName, _ := parameters[serverNameFieldID].(string)
	apiVersionStr, _ := parameters[apiVersionFieldID].(string)

	hostURL, err := url.Parse(hostURLStr)
	if err != nil {
		return nil, errors.New("pdns: invalid Host URL")
	}

	apiVersion, err := strconv.Atoi(apiVersionStr)
	if err != nil && apiVersionStr != "" {
		return nil, errors.New("pdns: invalid API version, must be an integer")
	}

	cfg := &pdns.Config{
		APIKey:             apiKey,
		Host:               hostURL,
		ServerName:         serverName,
		APIVersion:         apiVersion,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PoolingInterval,
		TTL:                dns.TTL,
	}

	return pdns.NewDNSProviderConfig(cfg)
}
