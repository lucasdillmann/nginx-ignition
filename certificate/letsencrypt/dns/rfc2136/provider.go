package rfc2136

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/rfc2136"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/ptr"
)

//nolint:gosec
const (
	nameserverFieldID    = "rfc2136Nameserver"
	tsigKeyFieldID       = "rfc2136TsigKey"
	tsigSecretFieldID    = "rfc2136TsigSecret"
	tsigAlgorithmFieldID = "rfc2136TsigAlgorithm"
)

type Provider struct{}

func (p *Provider) ID() string { return "RFC2136" }

func (p *Provider) Name() string { return "RFC2136" }

func (p *Provider) DynamicFields() []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          nameserverFieldID,
			Description: "DNS nameserver address",
			HelpText:    ptr.Of("host:port or host, defaults to port 53"),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          tsigKeyFieldID,
			Description: "TSIG key name",
			HelpText:    ptr.Of("Leave empty to disable TSIG authentication"),
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          tsigSecretFieldID,
			Description: "TSIG secret key",
			HelpText:    ptr.Of("Leave empty to disable TSIG authentication"),
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          tsigAlgorithmFieldID,
			Description: "TSIG algorithm",
			HelpText: ptr.Of(
				"e.g., hmac-sha256., defaults to hmac-sha1. Leave empty to disable TSIG authentication.",
			),
			Type: dynamicfields.SingleLineTextType,
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

	cfg := rfc2136.NewDefaultConfig()
	cfg.Nameserver = nameserver
	cfg.TSIGKey = tsigKey
	cfg.TSIGSecret = tsigSecret
	cfg.TSIGAlgorithm = tsigAlgorithm
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval
	cfg.TTL = dns.TTL

	return rfc2136.NewDNSProviderConfig(cfg)
}
