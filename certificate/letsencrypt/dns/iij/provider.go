package iij

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/iij"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

//nolint:gosec
const (
	accessKeyFieldID     = "iijAccessKey"
	secretKeyFieldID     = "iijSecretKey"
	doServiceCodeFieldID = "iijDoServiceCode"
)

type Provider struct{}

func (p *Provider) ID() string { return "IIJ" }

func (p *Provider) Name() string { return "IIJ" }

func (p *Provider) DynamicFields() []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          accessKeyFieldID,
			Description: "IIJ access key",
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          secretKeyFieldID,
			Description: "IIJ secret key",
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          doServiceCodeFieldID,
			Description: "IIJ service code",
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
	doServiceCode, _ := parameters[doServiceCodeFieldID].(string)

	cfg := &iij.Config{
		AccessKey:          accessKey,
		SecretKey:          secretKey,
		DoServiceCode:      doServiceCode,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PollingInterval,
		TTL:                dns.TTL,
	}

	return iij.NewDNSProviderConfig(cfg)
}
