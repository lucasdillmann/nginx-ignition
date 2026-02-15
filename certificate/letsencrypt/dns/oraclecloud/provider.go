package oraclecloud

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/oraclecloud"
	"github.com/nrdcg/oci-go-sdk/common/v1065"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

//nolint:gosec
const (
	compartmentOCIDFieldID   = "oracleCloudCompartmentOCID"
	regionFieldID            = "oracleCloudRegion"
	tenancyOCIDFieldID       = "oracleCloudTenancyOCID"
	userOCIDFieldID          = "oracleCloudUserOCID"
	pubKeyFingerprintFieldID = "oracleCloudPublicKeyFingerprint"
	privateKeyFieldID        = "oracleCloudPrivateKey"
	privateKeyPassFieldID    = "oracleCloudPrivateKeyPassphrase"
)

type Provider struct{}

func (p *Provider) ID() string { return "ORACLE_CLOUD" }

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateLetsencryptDnsOraclecloudName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID: compartmentOCIDFieldID,
			Description: i18n.M(
				ctx,
				i18n.K.CertificateLetsencryptDnsOraclecloudCompartmentOcid,
			),
			Required: true,
			Type:     dynamicfields.SingleLineTextType,
		},
		{
			ID:          regionFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsOraclecloudRegion),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          tenancyOCIDFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsOraclecloudTenancyOcid),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          userOCIDFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsOraclecloudUserOcid),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID: pubKeyFingerprintFieldID,
			Description: i18n.M(
				ctx,
				i18n.K.CertificateLetsencryptDnsOraclecloudPublicKeyFingerprint,
			),
			Required: true,
			Type:     dynamicfields.SingleLineTextType,
		},
		{
			ID:          privateKeyFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsOraclecloudPrivateKey),
			HelpText: i18n.M(
				ctx,
				i18n.K.CertificateLetsencryptDnsOraclecloudPrivateKeyHelp,
			),
			Required:  true,
			Sensitive: true,
			Type:      dynamicfields.MultiLineTextType,
		},
		{
			ID: privateKeyPassFieldID,
			Description: i18n.M(
				ctx,
				i18n.K.CertificateLetsencryptDnsOraclecloudPrivateKeyPassphrase,
			),
			Sensitive: true,
			Type:      dynamicfields.SingleLineTextType,
		},
	})
}

func (p *Provider) ChallengeProvider(
	_ context.Context,
	_ []string,
	parameters map[string]any,
) (challenge.Provider, error) {
	compartment, _ := parameters[compartmentOCIDFieldID].(string)
	region, _ := parameters[regionFieldID].(string)
	tenant, _ := parameters[tenancyOCIDFieldID].(string)
	user, _ := parameters[userOCIDFieldID].(string)
	finger, _ := parameters[pubKeyFingerprintFieldID].(string)
	privateKey, _ := parameters[privateKeyFieldID].(string)
	privateKeyPassword, _ := parameters[privateKeyPassFieldID].(string)

	cfg := oraclecloud.NewDefaultConfig()
	cfg.CompartmentID = compartment
	cfg.OCIConfigProvider = common.NewRawConfigurationProvider(
		tenant,
		user,
		region,
		finger,
		privateKey,
		new(privateKeyPassword),
	)
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval
	cfg.TTL = dns.TTL

	return oraclecloud.NewDNSProviderConfig(cfg)
}
