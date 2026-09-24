package exec

import (
	"context"

	"github.com/go-acme/lego/v5/challenge"
	"github.com/go-acme/lego/v5/providers/dns/exec"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/dynamicfields"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/i18n"
	dns2 "github.com/lucasdillmann/nginx-ignition/internal/certificate/domain/letsencrypt/dns"
)

const (
	programFieldID = "execProgram"
	modeFieldID    = "execMode"
)

type Provider struct{}

func (p *Provider) ID() string { return "EXEC_PROGRAM" }

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateLetsencryptDnsExecName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dns2.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          programFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsExecProgramPath),
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          modeFieldID,
			Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsExecRawMode),
			Type:        dynamicfields.BooleanType,
		},
	})
}

func (p *Provider) ChallengeProvider(
	_ context.Context,
	_ []string,
	parameters map[string]any,
) (challenge.Provider, error) {
	program, _ := parameters[programFieldID].(string)
	rawMode, _ := parameters[modeFieldID].(bool)

	mode := ""
	if rawMode {
		mode = "RAW"
	}

	cfg := exec.NewDefaultConfig()
	cfg.Program = program
	cfg.Mode = mode
	cfg.PropagationTimeout = dns2.PropagationTimeout
	cfg.PollingInterval = dns2.PollingInterval

	return exec.NewDNSProviderConfig(cfg)
}
