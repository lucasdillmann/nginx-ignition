package hurricane

import (
	"context"
	"strings"

	"github.com/go-acme/lego/v5/challenge"
	"github.com/go-acme/lego/v5/providers/dns/hurricane"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/coreerror"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/dynamicfields"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/i18n"
	dns2 "github.com/lucasdillmann/nginx-ignition/internal/certificate/domain/letsencrypt/dns"
)

const (
	tokensFieldID = "hurricaneTokens"
)

type Provider struct{}

func (p *Provider) ID() string { return "HURRICANE" }

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateLetsencryptDnsHurricaneName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns2.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          tokensFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsHurricaneTokens),
			HelpText:    i18n.M(ctx, i18n.K.CertificateLetsencryptDnsHurricaneTokensHelp),
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
	})
}

func (p *Provider) ChallengeProvider(
	ctx context.Context,
	_ []string,
	parameters map[string]any,
) (challenge.Provider, error) {
	tokensStr, _ := parameters[tokensFieldID].(string)

	credentials, err := parseTokens(ctx, tokensStr)
	if err != nil {
		return nil, err
	}

	cfg := hurricane.NewDefaultConfig()
	cfg.Credentials = credentials
	cfg.PropagationTimeout = dns2.PropagationTimeout
	cfg.PollingInterval = dns2.PollingInterval
	cfg.SequenceInterval = dns2.SequenceInterval

	return hurricane.NewDNSProviderConfig(cfg)
}

func parseTokens(ctx context.Context, tokensStr string) (map[string]string, error) {
	credentials := make(map[string]string)
	pairs := strings.Split(tokensStr, ",")

	for _, pair := range pairs {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) != 2 {
			return nil, coreerror.New(
				i18n.M(
					ctx,
					i18n.K.CertificateLetsencryptDnsHurricaneErrorHurricaneInvalidTokenFormat,
				),
				true,
			)
		}

		credentials[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	}

	return credentials, nil
}
