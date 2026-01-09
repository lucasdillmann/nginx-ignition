package scaleway

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/scaleway"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

//nolint:gosec
const (
	accessKeyFieldID = "scalewayAccessKey"
	secretKeyFieldID = "scalewaySecretKey"
	projectIDFieldID = "scalewayProjectId"
)

type Provider struct{}

func (p *Provider) ID() string { return "SCALEWAY" }

func (p *Provider) Name() string { return "Scaleway" }

func (p *Provider) DynamicFields() []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          accessKeyFieldID,
			Description: "Scaleway access key",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          secretKeyFieldID,
			Description: "Scaleway secret key",
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          projectIDFieldID,
			Description: "Scaleway project ID",
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
	accessKey, _ := parameters[accessKeyFieldID].(string)
	secretKey, _ := parameters[secretKeyFieldID].(string)
	projectID, _ := parameters[projectIDFieldID].(string)

	cfg := scaleway.NewDefaultConfig()
	cfg.AccessKey = accessKey
	cfg.Token = secretKey
	cfg.ProjectID = projectID
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval
	cfg.TTL = dns.TTL

	return scaleway.NewDNSProviderConfig(cfg)
}
