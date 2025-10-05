package exec

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/exec"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

const (
	programFieldID = "execProgram"
	modeFieldID    = "execMode"
)

type Provider struct{}

func (p *Provider) ID() string { return "EXEC_PROGRAM" }

func (p *Provider) Name() string { return "Custom script or program" }

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamic_fields.DynamicField{
		{
			ID:          programFieldID,
			Description: "Path to the program/script",
			Required:    true,
			Type:        dynamic_fields.SingleLineTextType,
		},
		{
			ID:          modeFieldID,
			Description: "Raw mode",
			Type:        dynamic_fields.BooleanType,
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
		PollingInterval:    dns.PoolingInterval,
	}

	return exec.NewDNSProviderConfig(cfg)
}
