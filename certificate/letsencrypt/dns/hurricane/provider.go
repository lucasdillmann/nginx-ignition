package hurricane

import (
	"context"
	"errors"
	"strings"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/hurricane"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/ptr"
)

const (
	tokensFieldID = "hurricaneTokens"
)

type Provider struct{}

func (p *Provider) ID() string { return "HURRICANE" }

func (p *Provider) Name() string { return "Hurricane Electric" }

func (p *Provider) DynamicFields() []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          tokensFieldID,
			Description: "Hurricane Electric tokens",
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
	tokensStr, _ := parameters[tokensFieldID].(string)

	credentials, err := parseTokens(tokensStr)
	if err != nil {
		return nil, err
	}

	cfg := hurricane.NewDefaultConfig()
	cfg.Credentials = credentials
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval
	cfg.SequenceInterval = dns.SequenceInterval

	return hurricane.NewDNSProviderConfig(cfg)
}

func parseTokens(tokensStr string) (map[string]string, error) {
	credentials := make(map[string]string)
	pairs := strings.Split(tokensStr, ",")

	for _, pair := range pairs {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) != 2 {
			return nil, errors.New("hurricane: invalid token format, expected key=value")
		}

		credentials[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	}

	return credentials, nil
}
