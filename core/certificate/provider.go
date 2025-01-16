package certificate

import "dillmann.com.br/nginx-ignition/core/common/dynamic_fields"

type Provider interface {
	ID() string
	Name() string
	DynamicFields() []*dynamic_fields.DynamicField
	Priority() int
	Issue(request *IssueRequest) (*bool, *Certificate, error)
	Renew(certificate *Certificate) (*bool, error)
}
