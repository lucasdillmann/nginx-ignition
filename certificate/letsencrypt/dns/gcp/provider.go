package gcp

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/gcloud"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

const (
	privateKeyID = "googleCloudPrivateKey"
)

type Provider struct{}

func (p *Provider) ID() string {
	return "GOOGLE_CLOUD"
}

func (p *Provider) Name() string {
	return "Google Cloud DNS"
}

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamic_fields.DynamicField{
		{
			ID:          privateKeyID,
			Description: "Service account private key JSON",
			Required:    true,
			Sensitive:   true,
			Type:        dynamic_fields.MultiLineTextType,
		},
	})
}

func (p *Provider) ChallengeProvider(
	_ context.Context,
	_ []string,
	parameters map[string]any,
) (challenge.Provider, error) {
	privateKey, _ := parameters[privateKeyID].(string)
	privateKeyBytes := []byte(privateKey)

	return gcloud.NewDNSProviderServiceAccountKey(privateKeyBytes)
}
