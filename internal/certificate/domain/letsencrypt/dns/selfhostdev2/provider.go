package selfhostdev2

import (
	"context"
	"strconv"
	"strings"

	"github.com/go-acme/lego/v5/challenge"
	"github.com/go-acme/lego/v5/providers/dns/selfhostdev2"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/coreerror"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/dynamicfields"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/i18n"
	"github.com/lucasdillmann/nginx-ignition/internal/certificate/domain/letsencrypt/dns"
)

//nolint:gosec
const (
	apiKeysFieldID        = "selfHostDeV2ApiKeys"
	recordsMappingFieldID = "selfHostDeV2RecordsMapping"
)

type Provider struct{}

func (p *Provider) ID() string {
	return "SELFHOSTDE_V2"
}

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateLetsencryptDnsSelfhostdev2Name)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          apiKeysFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsSelfhostdev2ApiKeys),
			HelpText:    i18n.M(ctx, i18n.K.CertificateLetsencryptDnsSelfhostdev2ApiKeysHelp),
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          recordsMappingFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsSelfhostdev2RecordsMapping),
			HelpText: i18n.M(
				ctx,
				i18n.K.CertificateLetsencryptDnsSelfhostdev2RecordsMappingHelp,
			),
			Required: true,
			Type:     dynamicfields.SingleLineTextType,
		},
	})
}

func (p *Provider) ChallengeProvider(
	ctx context.Context,
	_ []string,
	parameters map[string]any,
) (challenge.Provider, error) {
	apiKeys, _ := parameters[apiKeysFieldID].(string)
	recordsMapping, _ := parameters[recordsMappingFieldID].(string)

	credentials, err := parseAPIKeys(ctx, apiKeys)
	if err != nil {
		return nil, err
	}

	mapping, err := parseRecordsMapping(ctx, recordsMapping)
	if err != nil {
		return nil, err
	}

	cfg := selfhostdev2.NewDefaultConfig()
	cfg.Credentials = credentials
	cfg.RecordsMapping = mapping
	cfg.TTL = dns.TTL
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval

	return selfhostdev2.NewDNSProviderConfig(cfg)
}

func parseAPIKeys(ctx context.Context, raw string) (map[string]string, error) {
	credentials := make(map[string]string)

	for _, pair := range strings.Split(strings.TrimSpace(raw), ",") {
		parts := strings.Split(pair, ":")
		if len(parts) != 2 ||
			strings.TrimSpace(parts[0]) == "" ||
			strings.TrimSpace(parts[1]) == "" {
			return nil, coreerror.New(
				i18n.M(ctx, i18n.K.CertificateLetsencryptDnsSelfhostdev2InvalidApiKeysFormat),
				true,
			)
		}

		credentials[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	}

	return credentials, nil
}

func parseRecordsMapping(
	ctx context.Context,
	raw string,
) (map[string]*selfhostdev2.Seq, error) {
	mapping := make(map[string]*selfhostdev2.Seq)

	for _, pair := range strings.Split(strings.TrimSpace(raw), ",") {
		parts := strings.Split(pair, ":")
		domain := strings.TrimSpace(parts[0])

		if len(parts) < 2 || len(parts) > 3 || domain == "" {
			return nil, invalidRecordsMappingError(ctx)
		}

		ids := make([]int64, 0, len(parts)-1)

		for _, rawID := range parts[1:] {
			id, err := strconv.ParseInt(strings.TrimSpace(rawID), 10, 64)
			if err != nil {
				return nil, invalidRecordsMappingError(ctx)
			}

			ids = append(ids, id)
		}

		mapping[domain] = selfhostdev2.NewSeq(ids...)
	}

	return mapping, nil
}

func invalidRecordsMappingError(ctx context.Context) error {
	return coreerror.New(
		i18n.M(ctx, i18n.K.CertificateLetsencryptDnsSelfhostdev2InvalidRecordsMappingFormat),
		true,
	)
}
