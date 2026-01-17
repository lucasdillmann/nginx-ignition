package iijdpf

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/iijdpf"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

const (
	tokenFieldID       = "iijDpfToken"
	serviceCodeFieldID = "iijDpfServiceCode"
	endpointFieldID    = "iijDpfEndpoint"
)

type Provider struct{}

func (p *Provider) ID() string { return "IIJ_DPF" }

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsIijdpfName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          tokenFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsIijdpfApiToken),
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          serviceCodeFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsIijdpfServiceCode),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          endpointFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateCommonLetsEncryptDnsIijdpfApiEndpoint),
			Type:        dynamicfields.SingleLineTextType,
		},
	})
}

func (p *Provider) ChallengeProvider(
	_ context.Context,
	_ []string,
	parameters map[string]any,
) (challenge.Provider, error) {
	token, _ := parameters[tokenFieldID].(string)
	serviceCode, _ := parameters[serviceCodeFieldID].(string)
	endpoint, _ := parameters[endpointFieldID].(string)

	cfg := iijdpf.NewDefaultConfig()
	cfg.Token = token
	cfg.ServiceCode = serviceCode
	cfg.Endpoint = endpoint
	cfg.TTL = dns.TTL
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval

	return iijdpf.NewDNSProviderConfig(cfg)
}
