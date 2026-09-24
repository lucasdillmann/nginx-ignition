package certificate

import (
	"context"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/dynamicfields"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/i18n"
)

type Provider interface {
	ID() string
	Name(ctx context.Context) *i18n.Message
	DynamicFields(ctx context.Context) []dynamicfields.DynamicField
	Priority() int
	Issue(ctx context.Context, request *IssueRequest) (*Certificate, error)
	Renew(ctx context.Context, certificate *Certificate) (*Certificate, error)
	IsDueToRenew(ctx context.Context, certificate *Certificate) (bool, error)
}
