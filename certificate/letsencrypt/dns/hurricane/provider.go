package hurricane

import (
	"context"
	"errors"
	"strings"

	"github.com/aws/smithy-go/ptr"
	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/hurricane"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

const (
	tokensFieldID = "hurricaneTokens"
)

type Provider struct{}

func (p *Provider) ID() string { return "HURRICANE" }

func (p *Provider) Name() string { return "Hurricane Electric" }

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamic_fields.DynamicField{
		{
			ID:          tokensFieldID,
			Description: "Hurricane Electric tokens",
			HelpText:    ptr.String("comma-separated key=value pairs, e.g.: username=YOUR_USERNAME,password=YOUR_PASSWORD"),
			Required:    true,
			Sensitive:   true,
			Type:        dynamic_fields.SingleLineTextType,
		},
	})
}

func (p *Provider) ChallengeProvider(_ context.Context, _ []string, parameters map[string]any) (challenge.Provider, error) {
	tokensStr, _ := parameters[tokensFieldID].(string)

	credentials, err := parseTokens(tokensStr)
	if err != nil {
		return nil, err
	}

	cfg := &hurricane.Config{
		Credentials:        credentials,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PoolingInterval,
		SequenceInterval:   dns.PropagationTimeout,
	}

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
