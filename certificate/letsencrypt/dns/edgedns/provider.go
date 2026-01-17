package edgedns

import (
	"context"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/edgegrid"
	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/edgedns"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

//nolint:gosec
const (
	hostFieldID         = "edgeDnsHost"
	clientTokenFieldID  = "edgeDnsClientToken"
	clientSecretFieldID = "edgeDnsClientSecret"
	accessTokenFieldID  = "edgeDnsAccessToken"
	accountKeyFieldID   = "edgeDnsAccountKey"
)

type Provider struct{}

func (p *Provider) ID() string { return "EDGEDNS" }

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsEdgednsName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          hostFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsEdgednsEdgegridHost),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID: clientTokenFieldID,
			Description: i18n.M(
				ctx,
				i18n.K.CertificateCommonLetsEncryptDnsEdgednsEdgegridClientToken,
			),
			Required:  true,
			Sensitive: true,
			Type:      dynamicfields.SingleLineTextType,
		},
		{
			ID: clientSecretFieldID,
			Description: i18n.M(
				ctx,
				i18n.K.CertificateCommonLetsEncryptDnsEdgednsEdgegridClientSecret,
			),
			Required:  true,
			Sensitive: true,
			Type:      dynamicfields.SingleLineTextType,
		},
		{
			ID: accessTokenFieldID,
			Description: i18n.M(
				ctx,
				i18n.K.CertificateCommonLetsEncryptDnsEdgednsEdgegridAccessToken,
			),
			Required:  true,
			Sensitive: true,
			Type:      dynamicfields.SingleLineTextType,
		},
		{
			ID: accountKeyFieldID,
			Description: i18n.M(
				ctx,
				i18n.K.CertificateCommonLetsEncryptDnsEdgednsEdgegridAccountKey,
			),
			Required: true,
			Type:     dynamicfields.SingleLineTextType,
		},
	})
}

func (p *Provider) ChallengeProvider(
	_ context.Context,
	_ []string,
	parameters map[string]any,
) (challenge.Provider, error) {
	host, _ := parameters[hostFieldID].(string)
	clientToken, _ := parameters[clientTokenFieldID].(string)
	clientSecret, _ := parameters[clientSecretFieldID].(string)
	accessToken, _ := parameters[accessTokenFieldID].(string)
	accountKey, _ := parameters[accountKeyFieldID].(string)

	cfg := edgedns.NewDefaultConfig()
	cfg.Config = &edgegrid.Config{
		Host:         host,
		ClientToken:  clientToken,
		ClientSecret: clientSecret,
		AccessToken:  accessToken,
		AccountKey:   accountKey,
	}
	cfg.TTL = dns.TTL
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval

	return edgedns.NewDNSProviderConfig(cfg)
}
