package yandexcloud

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/yandexcloud"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

const (
	iamTokenFieldID = "yandexcloudIamToken"
	folderIDFieldID = "yandexcloudFolderId"
)

type Provider struct{}

func (p *Provider) ID() string {
	return "YANDEXCLOUD"
}

func (p *Provider) Name() string {
	return "Yandex Cloud"
}

func (p *Provider) DynamicFields() []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          iamTokenFieldID,
			Description: "Yandex Cloud IAM token",
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          folderIDFieldID,
			Description: "Yandex Cloud folder ID",
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

	cfg := &yandexcloud.Config{
		IamToken:           iamToken,
		FolderID:           folderID,
		TTL:                dns.TTL,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PollingInterval,
	}

	return yandexcloud.NewDNSProviderConfig(cfg)
}
