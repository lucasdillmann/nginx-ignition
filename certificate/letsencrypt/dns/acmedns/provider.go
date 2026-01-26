package acmedns

import (
	"context"
	"strings"

	"github.com/go-acme/lego/v4/challenge"
	legoacmedns "github.com/go-acme/lego/v4/providers/dns/acmedns"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

const (
	apiBaseFieldID        = "acmednsApiBase"
	allowListFieldID      = "acmednsAllowList"
	storagePathFieldID    = "acmednsStoragePath"
	storageBaseURLFieldID = "acmednsStorageBaseURL"
)

type Provider struct{}

func (p *Provider) ID() string { return "ACME_DNS" }

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateLetsencryptDnsAcmednsName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          apiBaseFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsAcmednsApiBase),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          allowListFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsAcmednsAllowList),
			HelpText:    i18n.M(ctx, i18n.K.CertificateLetsencryptDnsAcmednsAllowListHelp),
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          storagePathFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsAcmednsStoragePath),
			HelpText:    i18n.M(ctx, i18n.K.CertificateLetsencryptDnsAcmednsStoragePathHelp),
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          storageBaseURLFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsAcmednsStorageBaseUrl),
			HelpText: i18n.M(
				ctx,
				i18n.K.CertificateLetsencryptDnsAcmednsStorageBaseUrlHelp,
			),
			Type: dynamicfields.SingleLineTextType,
		},
	})
}

func (p *Provider) ChallengeProvider(
	_ context.Context,
	_ []string,
	parameters map[string]any,
) (challenge.Provider, error) {
	apiBase, _ := parameters[apiBaseFieldID].(string)
	allowListStr, _ := parameters[allowListFieldID].(string)
	storagePath, _ := parameters[storagePathFieldID].(string)
	storageBaseURL, _ := parameters[storageBaseURLFieldID].(string)

	cfg := legoacmedns.NewDefaultConfig()
	cfg.APIBase = apiBase
	cfg.StoragePath = storagePath
	cfg.StorageBaseURL = storageBaseURL

	if allowListStr != "" {
		list := make([]string, 0)
		for _, raw := range strings.Split(allowListStr, ",") {
			trimmedValue := strings.TrimSpace(raw)
			if trimmedValue != "" {
				list = append(list, trimmedValue)
			}
		}

		cfg.AllowList = list
	}

	return legoacmedns.NewDNSProviderConfig(cfg)
}
