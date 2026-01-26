package httpreq

import (
	"context"
	"net/url"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/httpreq"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/coreerror"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

const (
	endpointFieldID = "httpReqEndpoint"
	usernameFieldID = "httpReqUsername"
	passwordFieldID = "httpReqPassword"
	modeFieldID     = "httpReqMode"
)

type Provider struct{}

func (p *Provider) ID() string { return "HTTP_REQUEST" }

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateLetsencryptDnsHttpreqName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          endpointFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsHttpreqEndpointUrl),
			Required:    true,
			Type:        dynamicfields.URLType,
		},
		{
			ID:          usernameFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsHttpreqUsername),
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          passwordFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsHttpreqPassword),
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          modeFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsHttpreqRawMode),
			Type:        dynamicfields.BooleanType,
		},
	})
}

func (p *Provider) ChallengeProvider(
	ctx context.Context,
	_ []string,
	parameters map[string]any,
) (challenge.Provider, error) {
	endpointStr, _ := parameters[endpointFieldID].(string)
	username, _ := parameters[usernameFieldID].(string)
	password, _ := parameters[passwordFieldID].(string)
	rawMode, _ := parameters[modeFieldID].(bool)

	mode := ""
	if rawMode {
		mode = "RAW"
	}

	endpoint, err := url.Parse(endpointStr)
	if err != nil {
		return nil, coreerror.New(i18n.M(ctx, i18n.K.CommonInvalidUrl), true)
	}

	cfg := httpreq.NewDefaultConfig()
	cfg.Endpoint = endpoint
	cfg.Mode = mode
	cfg.Username = username
	cfg.Password = password
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval

	return httpreq.NewDNSProviderConfig(cfg)
}
