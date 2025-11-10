package nicmanager

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/nicmanager"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

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

func (p *Provider) Name() string { return "Nicmanager" }

func (p *Provider) DynamicFields() []*dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          loginFieldID,
			Description: "Nicmanager API login",
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          passwordFieldID,
			Description: "Nicmanager API password",
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          emailFieldID,
			Description: "Nicmanager API email",
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          usernameFieldID,
			Description: "Nicmanager API username",
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          otpSecretFieldID,
			Description: "Nicmanager OTP secret",
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          modeFieldID,
			Description: "Nicmanager API mode",
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

	cfg := &nicmanager.Config{
		Login:              login,
		Username:           username,
		Email:              email,
		Password:           password,
		OTPSecret:          otpSecret,
		Mode:               mode,
		TTL:                dns.TTL,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PollingInterval,
	}

	return nicmanager.NewDNSProviderConfig(cfg)
}
