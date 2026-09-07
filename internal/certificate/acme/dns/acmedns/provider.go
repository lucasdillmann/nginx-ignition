package acmedns

import (
	"context"
	"strings"

	"github.com/go-acme/lego/v5/challenge"
	legoacmedns "github.com/go-acme/lego/v5/providers/dns/acmedns"

	"github.com/lucasdillmann/nginx-ignition/internal/certificate/acme/dns"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/dynamicfields"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/i18n"
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
	return i18n.M(ctx, i18n.K.CertificateAcmeDnsAcmednsName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          apiBaseFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateAcmeDnsAcmednsApiBase),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          allowListFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateAcmeDnsAcmednsAllowList),
			HelpText:    i18n.M(ctx, i18n.K.CertificateAcmeDnsAcmednsAllowListHelp),
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          storagePathFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateAcmeDnsAcmednsStoragePath),
			HelpText:    i18n.M(ctx, i18n.K.CertificateAcmeDnsAcmednsStoragePathHelp),
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          storageBaseURLFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateAcmeDnsAcmednsStorageBaseUrl),
			HelpText: i18n.M(
				ctx,
				i18n.K.CertificateAcmeDnsAcmednsStorageBaseUrlHelp,
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
