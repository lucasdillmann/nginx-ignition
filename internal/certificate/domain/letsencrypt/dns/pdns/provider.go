package pdns

import (
	"context"
	"net/url"
	"strconv"

	"github.com/go-acme/lego/v5/challenge"
	"github.com/go-acme/lego/v5/providers/dns/pdns"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/coreerror"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/dynamicfields"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/i18n"
	dns2 "github.com/lucasdillmann/nginx-ignition/internal/certificate/domain/letsencrypt/dns"
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
	return i18n.M(ctx, i18n.K.CertificateLetsencryptDnsPdnsName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns2.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          apiKeyFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsPdnsApiKey),
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          hostURLFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsPdnsHostUrl),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          serverNameFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsPdnsServerName),
			HelpText:    i18n.M(ctx, i18n.K.CertificateLetsencryptDnsPdnsServerNameHelp),
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          apiVersionFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsPdnsApiVersion),
			HelpText:    i18n.M(ctx, i18n.K.CertificateLetsencryptDnsPdnsApiVersionHelp),
			Type:        dynamicfields.SingleLineTextType,
		},
	})
}

func (p *Provider) ChallengeProvider(
	ctx context.Context,
	_ []string,
	parameters map[string]any,
) (challenge.Provider, error) {
	apiKey, _ := parameters[apiKeyFieldID].(string)
	hostURLStr, _ := parameters[hostURLFieldID].(string)
	serverName, _ := parameters[serverNameFieldID].(string)
	apiVersionStr, _ := parameters[apiVersionFieldID].(string)

	hostURL, err := url.Parse(hostURLStr)
	if err != nil {
		return nil, coreerror.New(
			i18n.M(ctx, i18n.K.CertificateLetsencryptDnsPdnsErrorPdnsInvalidHostUrl),
			true,
		)
	}

	apiVersion, err := strconv.Atoi(apiVersionStr)
	if err != nil && apiVersionStr != "" {
		return nil, coreerror.New(
			i18n.M(ctx, i18n.K.CertificateLetsencryptDnsPdnsErrorPdnsInvalidApiVersion),
			true,
		)
	}

	cfg := pdns.NewDefaultConfig()
	cfg.APIKey = apiKey
	cfg.Host = hostURL
	cfg.ServerName = serverName
	cfg.APIVersion = apiVersion
	cfg.PropagationTimeout = dns2.PropagationTimeout
	cfg.PollingInterval = dns2.PollingInterval
	cfg.TTL = dns2.TTL

	return pdns.NewDNSProviderConfig(cfg)
}
