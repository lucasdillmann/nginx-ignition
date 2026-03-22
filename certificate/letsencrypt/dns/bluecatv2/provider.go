package bluecatv2

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/bluecatv2"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

const (
	serverURLFieldID  = "blueCatV2ServerUrl"
	usernameFieldID   = "blueCatV2Username"
	passwordFieldID   = "blueCatV2Password"
	configNameFieldID = "blueCatV2ConfigName"
	viewNameFieldID   = "blueCatV2ViewName"
)

type Provider struct{}

func (p *Provider) ID() string {
	return "BLUE_CAT_V2"
}

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateLetsencryptDnsBluecatv2Name)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          serverURLFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsBluecatv2ServerUrl),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          usernameFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsBluecatv2Username),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          passwordFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsBluecatv2Password),
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          configNameFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsBluecatv2ConfigName),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          viewNameFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsBluecatv2ViewName),
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
	config := bluecatv2.NewDefaultConfig()
	config.ServerURL, _ = parameters[serverURLFieldID].(string)
	config.Username, _ = parameters[usernameFieldID].(string)
	config.Password, _ = parameters[passwordFieldID].(string)
	config.ConfigName, _ = parameters[configNameFieldID].(string)
	config.ViewName, _ = parameters[viewNameFieldID].(string)

	return bluecatv2.NewDNSProviderConfig(config)
}
