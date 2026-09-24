package httpreq

import (
	"context"
	"net/url"

	"github.com/go-acme/lego/v5/challenge"
	"github.com/go-acme/lego/v5/providers/dns/httpreq"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/coreerror"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/dynamicfields"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/i18n"
	dns2 "github.com/lucasdillmann/nginx-ignition/internal/certificate/domain/letsencrypt/dns"
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
	return dns2.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
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
	cfg.PropagationTimeout = dns2.PropagationTimeout
	cfg.PollingInterval = dns2.PollingInterval

	return httpreq.NewDNSProviderConfig(cfg)
}
