package resolver

import (
	"context"

	"dillmann.com.br/nginx-ignition/core/integration"
)

type Option struct {
	*integration.DriverOption
	urlResolver func(ctx context.Context, option *Option) (*string, *[]string, error)
}

func (o *Option) URL(ctx context.Context) (*string, *[]string, error) {
	return o.urlResolver(ctx, o)
}
