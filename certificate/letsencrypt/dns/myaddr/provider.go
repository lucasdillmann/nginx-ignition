package myaddr

import (
	"context"
	"errors"
	"strings"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/myaddr"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

//nolint:gosec
const (
	credentialsFieldID = "myAddrCredentials"
)

type Provider struct{}

func (p *Provider) ID() string { return "MYADDR" }

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsMyaddrName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID: credentialsFieldID,
			Description: i18n.M(
				ctx,
				i18n.K.CertificateCommonLetsEncryptDnsMyaddrPrivateKeysMapping,
			),
			HelpText: i18n.M(
				ctx,
				i18n.K.CertificateCommonLetsEncryptDnsMyaddrPrivateKeysMappingHelp,
			),
			Required:  true,
			Sensitive: true,
			Type:      dynamicfields.SingleLineTextType,
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

	cfg := myaddr.NewDefaultConfig()
	cfg.Credentials = credentials
	cfg.TTL = dns.TTL
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval
	cfg.SequenceInterval = dns.SequenceInterval

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
