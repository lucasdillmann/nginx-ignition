package integration

import (
	"context"

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
	Name() string
	Description() string
	ConfigurationFields() []*dynamicfields.DynamicField
	GetAvailableOptions(
		ctx context.Context,
		parameters map[string]any,
		pageNumber, pageSize int,
		searchTerms *string,
		tcpOnly bool,
	) (*pagination.Page[*DriverOption], error)
	GetAvailableOptionById(
		ctx context.Context,
		parameters map[string]any,
		id string,
	) (*DriverOption, error)
	GetOptionProxyURL(
		ctx context.Context,
		parameters map[string]any,
		id string,
	) (*string, error)
}

type DriverOption struct {
	ID       string
	Name     string
	Port     int
	Protocol Protocol
}
