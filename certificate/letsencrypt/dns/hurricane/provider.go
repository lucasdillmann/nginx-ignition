package hurricane

import (
	"context"
	"strings"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/hurricane"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/coreerror"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

const (
	tokensFieldID = "hurricaneTokens"
)

type Provider struct{}

func (p *Provider) ID() string { return "HURRICANE" }

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsHurricaneName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          tokensFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsHurricaneTokens),
			HelpText:    i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsHurricaneTokensHelp),
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
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval
	cfg.SequenceInterval = dns.SequenceInterval

	return hurricane.NewDNSProviderConfig(cfg)
}

func parseTokens(ctx context.Context, tokensStr string) (map[string]string, error) {
	credentials := make(map[string]string)
	pairs := strings.Split(tokensStr, ",")

	for _, pair := range pairs {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) != 2 {
			return nil, coreerror.New(
				i18n.M(ctx, i18n.K.CertificateErrorHurricaneInvalidTokenFormat),
				true,
			)
		}

		credentials[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	}

	return credentials, nil
}
