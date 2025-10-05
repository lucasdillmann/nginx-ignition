package acmedns

import (
	"context"
	"strings"

	"github.com/aws/smithy-go/ptr"
	"github.com/go-acme/lego/v4/challenge"
	legoacmedns "github.com/go-acme/lego/v4/providers/dns/acmedns"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
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

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamic_fields.DynamicField{
		{
			ID:          apiBaseFieldID,
			Description: "ACME-DNS API base (e.g., https://acmedns.example.com)",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          allowListFieldID,
			Description: "Comma-separated list of CIDR sources allowed to update (optional)",
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          storagePathFieldID,
			Description: "Local storage file path for ACME-DNS accounts",
			HelpText:    ptr.String("Mutually exclusive with storage base URL"),
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          storageBaseURLFieldID,
			Description: "Remote storage base URL for ACME-DNS accounts",
			HelpText:    ptr.String("Mutually exclusive with storage path"),
			Type:        dynamic_fields.SingleLineTextType,
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
			s := strings.TrimSpace(raw)
			if s != "" {
				list = append(list, s)
			}
		}
		cfg.AllowList = list
	}

	return legoacmedns.NewDNSProviderConfig(cfg)
}
