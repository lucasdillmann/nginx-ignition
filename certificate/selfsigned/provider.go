package selfsigned

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
	return "SELF_SIGNED"
}

func (p *Provider) Name() string {
	return "Self-signed certificate"
}

func (p *Provider) DynamicFields() []*dynamic_fields.DynamicField {
	return []*dynamic_fields.DynamicField{}
}

func (p *Provider) Priority() int {
	return 3
}

func (p *Provider) Issue(_ *certificate.IssueRequest) (*certificate.Certificate, error) {
	// TODO: Implement this
	return nil, core_error.New("not implemented yet", false)
}

func (p *Provider) Renew(_ *certificate.Certificate) (*certificate.Certificate, error) {
	// TODO: Implement this
	return nil, core_error.New("not implemented yet", false)
}
