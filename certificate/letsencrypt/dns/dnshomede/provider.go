package dnshomede

import (
	"context"
	"errors"
	"strings"

	"github.com/aws/smithy-go/ptr"
	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/dnshomede"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

const (
	credentialsFieldID = "dnsHomeDeCredentials"
)

type Provider struct{}

func (p *Provider) ID() string { return "DNSHOME_DE" }

func (p *Provider) Name() string { return "dnsHome.de" }

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamic_fields.DynamicField{
		{
			ID:          credentialsFieldID,
			Description: "dnsHome.de credentials",
			HelpText:    ptr.String("Comma-separated key=value pairs for authentication"),
			Required:    true,
			Sensitive:   true,
			Type:        dynamic_fields.SingleLineTextType,
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

	cfg := &dnshomede.Config{
		Credentials:        credentials,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PoolingInterval,
		SequenceInterval:   dns.SequenceInterval,
	}

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
