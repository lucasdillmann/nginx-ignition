package oraclecloud

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/oraclecloud"
	"github.com/nrdcg/oci-go-sdk/common/v1065"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/ptr"
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

func (p *Provider) Name() string { return "Oracle Cloud (OCI)" }

func (p *Provider) DynamicFields() []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          compartmentOCIDFieldID,
			Description: "OCI compartment OCID",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          regionFieldID,
			Description: "OCI region",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          tenancyOCIDFieldID,
			Description: "OCI tenancy OCID",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          userOCIDFieldID,
			Description: "OCI user OCID",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          pubKeyFingerprintFieldID,
			Description: "OCI public key fingerprint",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          privateKeyFieldID,
			Description: "OCI private key",
			HelpText:    ptr.Of("PEM contents"),
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.MultiLineTextType,
		},
		{
			ID:          privateKeyPassFieldID,
			Description: "OCI private key passphrase",
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
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

	cfg := &oraclecloud.Config{
		CompartmentID: compartment,
		OCIConfigProvider: common.NewRawConfigurationProvider(
			tenant,
			user,
			region,
			finger,
			privateKey,
			ptr.Of(privateKeyPassword),
		),
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PollingInterval,
		TTL:                dns.TTL,
	}

	return oraclecloud.NewDNSProviderConfig(cfg)
}
