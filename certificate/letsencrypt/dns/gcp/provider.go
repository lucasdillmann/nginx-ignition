package gcp

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/gcloud"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

const (
	privateKeyFieldID = "googleCloudPrivateKey"
)

type Provider struct{}

func (p *Provider) ID() string {
	return "GOOGLE_CLOUD"
}

func (p *Provider) Name() string {
	return "Google Cloud"
}

func (p *Provider) DynamicFields() []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          privateKeyFieldID,
			Description: "Service account private key JSON",
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.MultiLineTextType,
		},
	})
}

func (p *Provider) ChallengeProvider(
	_ context.Context,
	_ []string,
	parameters map[string]any,
) (challenge.Provider, error) {
	privateKey, _ := parameters[privateKeyFieldID].(string)
	privateKeyBytes := []byte(privateKey)

	return gcloud.NewDNSProviderServiceAccountKey(privateKeyBytes)
}
