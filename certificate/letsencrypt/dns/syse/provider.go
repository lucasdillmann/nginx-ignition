package syse

import (
	"context"
	"strings"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/syse"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

const (
	credentialsFieldID = "syseCredentials"
)

type Provider struct{}

func (p *Provider) ID() string { return "SYSE" }

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateLetsencryptDnsSyseName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          credentialsFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsSyseCredentials),
			HelpText:    i18n.M(ctx, i18n.K.CertificateLetsencryptDnsSyseCredentialsHelp),
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.MultiLineTextType,
		},
	})
}

func (p *Provider) ChallengeProvider(
	_ context.Context,
	_ []string,
	parameters map[string]any,
) (challenge.Provider, error) {
	raw, _ := parameters[credentialsFieldID].(string)

	credentials := make(map[string]string)
	for _, line := range strings.Split(raw, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			credentials[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}

	cfg := syse.NewDefaultConfig()
	cfg.Credentials = credentials
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval
	cfg.TTL = dns.TTL

	return syse.NewDNSProviderConfig(cfg)
}
