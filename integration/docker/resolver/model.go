package resolver

import (
	"context"

	"dillmann.com.br/nginx-ignition/core/integration"
)

type Resolver interface {
	ResolveOptions(ctx context.Context, tcpOnly bool, searchTerms *string) (*[]Option, error)
	ResolveOptionByID(ctx context.Context, optionId string) (*Option, error)
}

type Option struct {
	*integration.DriverOption
	urlResolver func(ctx context.Context, option *Option, publicUrl string) (*string, error)
}

func (o *Option) URL(ctx context.Context, publicUrl string) (*string, error) {
	return o.urlResolver(ctx, o, publicUrl)
}

const (
	hostQualifier      = "host"
	containerQualifier = "container"
	ingressQualifier   = "ingress"
)
