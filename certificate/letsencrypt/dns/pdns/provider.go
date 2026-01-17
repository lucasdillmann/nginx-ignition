package pdns

import (
	"context"
	"errors"
	"net/url"
	"strconv"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/pdns"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

//nolint:gosec
const (
	apiKeyFieldID     = "powerDnsApiKey"
	hostURLFieldID    = "powerDnsHostUtl"
	serverNameFieldID = "powerDnsServerName"
	apiVersionFieldID = "powerDnsApiVersion"
)

type Provider struct{}

func (p *Provider) ID() string { return "POWERDNS" }

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsPdnsName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          apiKeyFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsPdnsApiKey),
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          hostURLFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsPdnsHostUrl),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          serverNameFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsPdnsServerName),
			HelpText:    i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsPdnsServerNameHelp),
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          apiVersionFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsPdnsApiVersion),
			HelpText:    i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsPdnsApiVersionHelp),
			Type:        dynamicfields.SingleLineTextType,
		},
	})
}

func (p *Provider) ChallengeProvider(
	_ context.Context,
	_ []string,
	parameters map[string]any,
) (challenge.Provider, error) {
	apiKey, _ := parameters[apiKeyFieldID].(string)
	hostURLStr, _ := parameters[hostURLFieldID].(string)
	serverName, _ := parameters[serverNameFieldID].(string)
	apiVersionStr, _ := parameters[apiVersionFieldID].(string)

	hostURL, err := url.Parse(hostURLStr)
	if err != nil {
		return nil, errors.New("pdns: invalid Host URL")
	}

	apiVersion, err := strconv.Atoi(apiVersionStr)
	if err != nil && apiVersionStr != "" {
		return nil, errors.New("pdns: invalid API version, must be an integer")
	}

	cfg := pdns.NewDefaultConfig()
	cfg.APIKey = apiKey
	cfg.Host = hostURL
	cfg.ServerName = serverName
	cfg.APIVersion = apiVersion
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval
	cfg.TTL = dns.TTL

	return pdns.NewDNSProviderConfig(cfg)
}
