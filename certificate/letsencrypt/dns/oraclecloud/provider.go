package oraclecloud

import (
	"context"

	"github.com/aws/smithy-go/ptr"
	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/oraclecloud"
	"github.com/nrdcg/oci-go-sdk/common/v1065"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

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

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamic_fields.DynamicField{
		{
			ID:          compartmentOCIDFieldID,
			Description: "OCI compartment OCID",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          regionFieldID,
			Description: "OCI region",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          tenancyOCIDFieldID,
			Description: "OCI tenancy OCID",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          userOCIDFieldID,
			Description: "OCI user OCID",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          pubKeyFingerprintFieldID,
			Description: "OCI public key fingerprint",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          privateKeyFieldID,
			Description: "OCI private key",
			HelpText:    ptr.String("PEM contents"),
			Required:    true,
			Sensitive:   true,
			Type:        dynamic_fields.MultiLineTextType,
		},
		{
			ID:          privateKeyPassFieldID,
			Description: "OCI private key passphrase",
			Sensitive:   true,
			Type:        dynamic_fields.SingleLineTextType,
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
			ptr.String(privateKeyPassword),
		),
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PoolingInterval,
		TTL:                dns.TTL,
	}

	return oraclecloud.NewDNSProviderConfig(cfg)
}
