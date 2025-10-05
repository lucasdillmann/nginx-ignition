package rfc2136

import (
	"context"

	"github.com/aws/smithy-go/ptr"
	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/rfc2136"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

const (
	nameserverFieldID    = "rfc2136Nameserver"
	tsigKeyFieldID       = "rfc2136TsigKey"
	tsigSecretFieldID    = "rfc2136TsigSecret"
	tsigAlgorithmFieldID = "rfc2136TsigAlgorithm"
)

type Provider struct{}

func (p *Provider) ID() string { return "RFC2136" }

func (p *Provider) Name() string { return "RFC2136" }

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamic_fields.DynamicField{
		{
			ID:          nameserverFieldID,
			Description: "DNS nameserver address",
			HelpText:    ptr.String("host:port or host, defaults to port 53"),
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          tsigKeyFieldID,
			Description: "TSIG key name",
			HelpText:    ptr.String("Leave empty to disable TSIG authentication"),
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          tsigSecretFieldID,
			Description: "TSIG secret key",
			HelpText:    ptr.String("Leave empty to disable TSIG authentication"),
			Sensitive:   true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          tsigAlgorithmFieldID,
			Description: "TSIG algorithm",
			HelpText:    ptr.String("e.g., hmac-sha256., defaults to hmac-sha1. Leave empty to disable TSIG authentication."),
			Type:        dynamic_fields.SingleLineTextType,
		},
	})
}

func (p *Provider) ChallengeProvider(
	_ context.Context,
	_ []string,
	parameters map[string]any,
) (challenge.Provider, error) {
	nameserver, _ := parameters[nameserverFieldID].(string)
	tsigKey, _ := parameters[tsigKeyFieldID].(string)
	tsigSecret, _ := parameters[tsigSecretFieldID].(string)
	tsigAlgorithm, _ := parameters[tsigAlgorithmFieldID].(string)

	cfg := &rfc2136.Config{
		Nameserver:         nameserver,
		TSIGKey:            tsigKey,
		TSIGSecret:         tsigSecret,
		TSIGAlgorithm:      tsigAlgorithm,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PollingInterval,
		TTL:                dns.TTL,
	}

	return rfc2136.NewDNSProviderConfig(cfg)
}
