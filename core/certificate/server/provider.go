package server

import (
	"context"

	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

type Provider interface {
	ID() string
	Name() string
	DynamicFields() []*dynamicfields.DynamicField
	Priority() int
	Issue(ctx context.Context, request *IssueRequest) (*Certificate, error)
	Renew(ctx context.Context, certificate *Certificate) (*Certificate, error)
}
