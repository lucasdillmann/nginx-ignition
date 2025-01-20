package letsencrypt

import (
	"dillmann.com.br/nginx-ignition/core/certificate"
	"dillmann.com.br/nginx-ignition/core/common/core_error"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

type Provider struct{}

func New() *Provider {
	return &Provider{}
}

func (p *Provider) ID() string {
	return "LETS_ENCRYPT"
}

func (p *Provider) Name() string {
	return "Let's encrypt"
}

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return []*dynamic_fields.DynamicField{}
}

func (p *Provider) Priority() int {
	return 1
}

func (p *Provider) Issue(_ *certificate.IssueRequest) (*certificate.Certificate, error) {
	return nil, core_error.New("not implemented yet", false)
}

func (p *Provider) Renew(_ *certificate.Certificate) (*certificate.Certificate, error) {
	return nil, core_error.New("not implemented yet", false)
}
