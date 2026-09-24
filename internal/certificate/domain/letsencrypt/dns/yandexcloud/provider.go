package yandexcloud

import (
	"context"

	"github.com/go-acme/lego/v5/challenge"
	"github.com/go-acme/lego/v5/providers/dns/yandexcloud"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/dynamicfields"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/i18n"
	dns2 "github.com/lucasdillmann/nginx-ignition/internal/certificate/domain/letsencrypt/dns"
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
	return dns2.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
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
	cfg.TTL = dns2.TTL
	cfg.PropagationTimeout = dns2.PropagationTimeout
	cfg.PollingInterval = dns2.PollingInterval

	return yandexcloud.NewDNSProviderConfig(cfg)
}
