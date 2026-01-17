package integration

import (
	"context"

	"dillmann.com.br/nginx-ignition/core/common/i18n"

	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/pagination"
)

type Protocol string

const (
	TCPProtocol Protocol = "TCP"
	UDPProtocol Protocol = "UDP"
)

type Driver interface {
	ID() string
	Name(ctx context.Context) *i18n.Message
	Description(ctx context.Context) *i18n.Message
	ConfigurationFields(ctx context.Context) []dynamicfields.DynamicField
	GetAvailableOptions(
		ctx context.Context,
		parameters map[string]any,
		pageNumber, pageSize int,
		searchTerms *string,
		tcpOnly bool,
	) (*pagination.Page[DriverOption], error)
	GetAvailableOptionByID(
		ctx context.Context,
		parameters map[string]any,
		id string,
	) (*DriverOption, error)
	GetOptionProxyURL(
		ctx context.Context,
		parameters map[string]any,
		id string,
	) (*string, []string, error)
}

type DriverOption struct {
	Qualifier    *string
	ID           string
	Name         string
	Protocol     Protocol
	DNSResolvers []string
	Port         int
}
