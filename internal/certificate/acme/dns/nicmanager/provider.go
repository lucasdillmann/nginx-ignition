package nicmanager

import (
	"context"

	"github.com/go-acme/lego/v5/challenge"
	"github.com/go-acme/lego/v5/providers/dns/nicmanager"

	"github.com/lucasdillmann/nginx-ignition/internal/certificate/acme/dns"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/dynamicfields"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/i18n"
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
	return i18n.M(ctx, i18n.K.CertificateAcmeDnsNicmanagerName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          loginFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateAcmeDnsNicmanagerApiLogin),
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          passwordFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateAcmeDnsNicmanagerApiPassword),
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          emailFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateAcmeDnsNicmanagerApiEmail),
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          usernameFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateAcmeDnsNicmanagerApiUsername),
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          otpSecretFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateAcmeDnsNicmanagerOtpSecret),
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          modeFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateAcmeDnsNicmanagerApiMode),
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
