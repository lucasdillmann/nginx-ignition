package myaddr

import (
	"context"
	"strings"

	"github.com/go-acme/lego/v5/challenge"
	"github.com/go-acme/lego/v5/providers/dns/myaddr"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/coreerror"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/dynamicfields"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/i18n"
	dns2 "github.com/lucasdillmann/nginx-ignition/internal/certificate/domain/letsencrypt/dns"
)

//nolint:gosec
const (
	credentialsFieldID = "myAddrCredentials"
)

type Provider struct{}

func (p *Provider) ID() string { return "MYADDR" }

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateLetsencryptDnsMyaddrName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns2.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID: credentialsFieldID,
			Description: i18n.M(
				ctx,
				i18n.K.CertificateLetsencryptDnsMyaddrPrivateKeysMapping,
			),
			HelpText: i18n.M(
				ctx,
				i18n.K.CertificateLetsencryptDnsMyaddrPrivateKeysMappingHelp,
			),
			Required:  true,
			Sensitive: true,
			Type:      dynamicfields.SingleLineTextType,
		},
	})
}

func (p *Provider) ChallengeProvider(
	ctx context.Context,
	_ []string,
	parameters map[string]any,
) (challenge.Provider, error) {
	credentialsStr, _ := parameters[credentialsFieldID].(string)

	credentials, err := parseCredentials(ctx, credentialsStr)
	if err != nil {
		return nil, err
	}

	cfg := myaddr.NewDefaultConfig()
	cfg.Credentials = credentials
	cfg.TTL = dns2.TTL
	cfg.PropagationTimeout = dns2.PropagationTimeout
	cfg.PollingInterval = dns2.PollingInterval
	cfg.SequenceInterval = dns2.SequenceInterval

	return myaddr.NewDNSProviderConfig(cfg)
}

func parseCredentials(ctx context.Context, credentialsStr string) (map[string]string, error) {
	credentials := make(map[string]string)
	pairs := strings.Split(credentialsStr, ",")

	for _, pair := range pairs {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) != 2 {
			return nil, coreerror.New(
				i18n.M(
					ctx,
					i18n.K.CertificateLetsencryptDnsMyaddrErrorMyaddrInvalidCredentialsFormat,
				),
				true,
			)
		}

		credentials[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	}

	return credentials, nil
}
