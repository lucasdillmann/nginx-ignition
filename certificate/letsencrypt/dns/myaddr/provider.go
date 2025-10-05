package myaddr

import (
	"context"
	"errors"
	"strings"

	"github.com/aws/smithy-go/ptr"
	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/myaddr"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

const (
	credentialsFieldID = "myAddrCredentials"
)

type Provider struct{}

func (p *Provider) ID() string { return "MYADDR" }

func (p *Provider) Name() string { return "myaddr" }

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamic_fields.DynamicField{
		{
			ID:          credentialsFieldID,
			Description: "myaddr private keys mapping",
			HelpText:    ptr.String("Comma-separated key=value pairs"),
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

	cfg := &myaddr.Config{
		Credentials:        credentials,
		TTL:                dns.TTL,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PollingInterval,
		SequenceInterval:   dns.SequenceInterval,
	}

	return myaddr.NewDNSProviderConfig(cfg)
}

func parseCredentials(credentialsStr string) (map[string]string, error) {
	credentials := make(map[string]string)
	pairs := strings.Split(credentialsStr, ",")

	for _, pair := range pairs {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) != 2 {
			return nil, errors.New("myaddr: invalid credentials format, expected key=value")
		}

		credentials[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	}

	return credentials, nil
}
