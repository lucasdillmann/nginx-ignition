package yandexcloud

import (
	"context"

	"github.com/go-acme/lego/v5/challenge"
	"github.com/go-acme/lego/v5/providers/dns/yandexcloud"

	"nginx-ignition/internal/certificate/letsencrypt/dns"
	"nginx-ignition/internal/core/common/dynamicfields"
	"nginx-ignition/internal/core/common/i18n"
)

const (
	iamTokenFieldID = "yandexcloudIamToken"
	folderIDFieldID = "yandexcloudFolderId"
)

type Provider struct{}

func (p *Provider) ID() string {
	return "YANDEXCLOUD"
}

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateLetsencryptDnsYandexcloudName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          iamTokenFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsYandexcloudIamToken),
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          folderIDFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsYandexcloudFolderId),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
	})
}

func (p *Provider) ChallengeProvider(
	_ context.Context,
	_ []string,
	parameters map[string]any,
) (challenge.Provider, error) {
	iamToken, _ := parameters[iamTokenFieldID].(string)
	folderID, _ := parameters[folderIDFieldID].(string)

	cfg := yandexcloud.NewDefaultConfig()
	cfg.IamToken = iamToken
	cfg.FolderID = folderID
	cfg.TTL = dns.TTL
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval

	return yandexcloud.NewDNSProviderConfig(cfg)
}
