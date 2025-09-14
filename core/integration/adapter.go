package integration

import (
	"context"

	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
	"dillmann.com.br/nginx-ignition/core/common/pagination"
)

type Protocol string

const (
	TCPProtocol Protocol = "TCP"
	UDPProtocol Protocol = "UDP"
)

type Adapter interface {
	ID() string
	Name() string
	Priority() int
	Description() string
	ConfigurationFields() []*dynamic_fields.DynamicField
	GetAvailableOptions(
		ctx context.Context,
		parameters map[string]any,
		pageNumber, pageSize int,
		searchTerms *string,
		tcpOnly bool,
	) (*pagination.Page[*AdapterOption], error)
	GetAvailableOptionById(
		ctx context.Context,
		parameters map[string]any,
		id string,
	) (*AdapterOption, error)
	GetOptionProxyUrl(
		ctx context.Context,
		parameters map[string]any,
		id string,
	) (*string, error)
}

type AdapterOption struct {
	ID       string
	Name     string
	Port     int
	Protocol Protocol
}
