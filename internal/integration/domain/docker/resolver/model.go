package resolver

import (
	"context"

	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/integration"
)

type Option struct {
	urlResolver func(ctx context.Context, option *Option) (*string, []string, error)
	integration.DriverOption
	privatePort int
}

func (o *Option) URL(ctx context.Context) (*string, []string, error) {
	return o.urlResolver(ctx, o)
}
