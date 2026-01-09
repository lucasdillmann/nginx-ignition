package dnshomede

import (
	"context"
	"errors"
	"strings"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/dnshomede"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/ptr"
)

//nolint:gosec
const (
	credentialsFieldID = "dnsHomeDeCredentials"
)

type Provider struct{}

func (p *Provider) ID() string { return "DNSHOME_DE" }

func (p *Provider) Name() string { return "dnsHome.de" }

func (p *Provider) DynamicFields() []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          credentialsFieldID,
			Description: "dnsHome.de credentials",
			HelpText:    ptr.Of("Comma-separated key=value pairs"),
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
	credentialsStr, _ := parameters[credentialsFieldID].(string)

	credentials, err := parseCredentials(credentialsStr)
	if err != nil {
		return nil, err
	}

	cfg := dnshomede.NewDefaultConfig()
	cfg.Credentials = credentials
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval
	cfg.SequenceInterval = dns.SequenceInterval

	return dnshomede.NewDNSProviderConfig(cfg)
}

func parseCredentials(credentialsStr string) (map[string]string, error) {
	credentials := make(map[string]string)
	pairs := strings.Split(credentialsStr, ",")

	for _, pair := range pairs {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) != 2 {
			return nil, errors.New("dnshomede: invalid credentials format, expected key=value")
		}

		credentials[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	}

	return credentials, nil
}
