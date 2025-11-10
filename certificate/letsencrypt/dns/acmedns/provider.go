package acmedns

import (
	"context"
	"strings"

	"github.com/go-acme/lego/v4/challenge"
	legoacmedns "github.com/go-acme/lego/v4/providers/dns/acmedns"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/ptr"
)

const (
	apiBaseFieldID        = "acmednsApiBase"
	allowListFieldID      = "acmednsAllowList"
	storagePathFieldID    = "acmednsStoragePath"
	storageBaseURLFieldID = "acmednsStorageBaseURL"
)

type Provider struct{}

func (p *Provider) ID() string { return "ACME_DNS" }

func (p *Provider) Name() string { return "ACME-DNS" }

func (p *Provider) DynamicFields() []*dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          apiBaseFieldID,
			Description: "ACME-DNS API base",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          allowListFieldID,
			Description: "CIDR sources allowed to update",
			HelpText:    ptr.Of("Comma-separated key=value pairs"),
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          storagePathFieldID,
			Description: "Local storage file path for ACME-DNS accounts",
			HelpText:    ptr.Of("Mutually exclusive with storage base URL"),
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          storageBaseURLFieldID,
			Description: "Remote storage base URL for ACME-DNS accounts",
			HelpText:    ptr.Of("Mutually exclusive with storage path"),
			Type:        dynamicfields.SingleLineTextType,
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

	cfg := &legoacmedns.Config{
		APIBase:        apiBase,
		StoragePath:    storagePath,
		StorageBaseURL: storageBaseURL,
	}

	if allowListStr != "" {
		var list []string
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
