package exec

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/exec"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

const (
	programFieldID = "execProgram"
	modeFieldID    = "execMode"
)

type Provider struct{}

func (p *Provider) ID() string { return "EXEC_PROGRAM" }

func (p *Provider) Name() string { return "Custom script or program" }

func (p *Provider) DynamicFields() []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          programFieldID,
			Description: "Path to the program/script",
			Required:    true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          modeFieldID,
			Description: "Raw mode",
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

	cfg := &exec.Config{
		Program:            program,
		Mode:               mode,
		PropagationTimeout: dns.PropagationTimeout,
		PollingInterval:    dns.PollingInterval,
	}

	return exec.NewDNSProviderConfig(cfg)
}
