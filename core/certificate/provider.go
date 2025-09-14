package certificate

import (
	"context"

	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

type Provider interface {
	ID() string
	Name() string
	DynamicFields() []*dynamic_fields.DynamicField
	Priority() int
	Issue(ctx context.Context, request *IssueRequest) (*Certificate, error)
	Renew(ctx context.Context, certificate *Certificate) (*Certificate, error)
}
