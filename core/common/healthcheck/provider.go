package healthcheck

import (
	"context"
)

type Provider interface {
	ID() string
	Check(ctx context.Context) error
}
