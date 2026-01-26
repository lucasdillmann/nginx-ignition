package nicmanager

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/nicmanager"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

//nolint:gosec
const (
	loginFieldID     = "nicManagerLogin"
	usernameFieldID  = "nicManagerUsername"
	emailFieldID     = "nicManagerEmail"
	passwordFieldID  = "nicManagerPassword"
	otpSecretFieldID = "nicManagerOtpSecret"
	modeFieldID      = "nicManagerMode"
)

type Provider struct{}

func (p *Provider) ID() string { return "NICMANAGER" }

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateLetsencryptDnsNicmanagerName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          loginFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsNicmanagerApiLogin),
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          passwordFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsNicmanagerApiPassword),
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          emailFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsNicmanagerApiEmail),
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          usernameFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsNicmanagerApiUsername),
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          otpSecretFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsNicmanagerOtpSecret),
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          modeFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsNicmanagerApiMode),
			Type:        dynamicfields.SingleLineTextType,
		},
	})
}

func (p *Provider) ChallengeProvider(
	_ context.Context,
	_ []string,
	parameters map[string]any,
) (challenge.Provider, error) {
	login, _ := parameters[loginFieldID].(string)
	username, _ := parameters[usernameFieldID].(string)
	email, _ := parameters[emailFieldID].(string)
	password, _ := parameters[passwordFieldID].(string)
	otpSecret, _ := parameters[otpSecretFieldID].(string)
	mode, _ := parameters[modeFieldID].(string)

	cfg := nicmanager.NewDefaultConfig()
	cfg.Login = login
	cfg.Username = username
	cfg.Email = email
	cfg.Password = password
	cfg.OTPSecret = otpSecret
	cfg.Mode = mode
	cfg.TTL = dns.TTL
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval

	return nicmanager.NewDNSProviderConfig(cfg)
}
