package fields

import (
	"context"

	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

const (
	SocketConnectionMode = "SOCKET"
	TCPConnectionMode    = "TCP"

	ConnectionModeFieldID       = "connectionMode"
	SocketPathFieldID           = "socketPath"
	HostURLFieldID              = "hostUrl"
	SwarmModeFieldID            = "swarmMode"
	SwarmServiceMeshFieldID     = "swarmServiceMesh"
	SwarmDNSResolversFieldID    = "swarmDnsResolvers"
	UseContainerNameAsIDFieldID = "useContainerNameAsId"
	ProxyURLFieldID             = "proxyUrl"
)

func DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	connectionMode := dynamicfields.DynamicField{
		ID:           ConnectionModeFieldID,
		Description:  i18n.M(ctx, i18n.K.IntegrationDockerFieldsConnectionMode),
		Priority:     1,
		Required:     true,
		Type:         dynamicfields.EnumType,
		DefaultValue: SocketConnectionMode,
		EnumOptions: []dynamicfields.EnumOption{
			{
				ID:          SocketConnectionMode,
				Description: i18n.M(ctx, i18n.K.IntegrationDockerFieldsConnectionModeSocket),
			},
			{ID: TCPConnectionMode, Description: i18n.M(ctx, i18n.K.IntegrationDockerFieldsConnectionModeTcp)},
		},
	}

	socketPath := dynamicfields.DynamicField{
		ID:           SocketPathFieldID,
		Description:  i18n.M(ctx, i18n.K.IntegrationDockerFieldsSocketPath),
		Priority:     2,
		Required:     true,
		Type:         dynamicfields.SingleLineTextType,
		DefaultValue: "/var/run/docker.sock",
		Conditions: []dynamicfields.Condition{{
			ParentField: ConnectionModeFieldID,
			Value:       SocketConnectionMode,
		}},
	}

	hostURL := dynamicfields.DynamicField{
		ID:           HostURLFieldID,
		Description:  i18n.M(ctx, i18n.K.IntegrationDockerFieldsHostUrl),
		Priority:     3,
		Required:     true,
		Type:         dynamicfields.URLType,
		DefaultValue: "",
		HelpText:     i18n.M(ctx, i18n.K.IntegrationDockerFieldsHostUrlHelp),
		Conditions: []dynamicfields.Condition{{
			ParentField: ConnectionModeFieldID,
			Value:       TCPConnectionMode,
		}},
	}

	swarmMode := dynamicfields.DynamicField{
		ID:           SwarmModeFieldID,
		Description:  i18n.M(ctx, i18n.K.IntegrationDockerFieldsSwarmMode),
		Priority:     4,
		Required:     true,
		Type:         dynamicfields.BooleanType,
		DefaultValue: false,
		HelpText:     i18n.M(ctx, i18n.K.IntegrationDockerFieldsSwarmModeHelp),
	}

	swarmServiceMesh := dynamicfields.DynamicField{
		ID:           SwarmServiceMeshFieldID,
		Description:  i18n.M(ctx, i18n.K.IntegrationDockerFieldsSwarmServiceMesh),
		Priority:     5,
		Required:     true,
		Type:         dynamicfields.BooleanType,
		DefaultValue: false,
		HelpText:     i18n.M(ctx, i18n.K.IntegrationDockerFieldsSwarmServiceMeshHelp),
		Conditions: []dynamicfields.Condition{{
			ParentField: SwarmModeFieldID,
			Value:       true,
		}},
	}

	swarmDNSResolvers := dynamicfields.DynamicField{
		ID:           SwarmDNSResolversFieldID,
		Description:  i18n.M(ctx, i18n.K.IntegrationDockerFieldsSwarmDnsResolvers),
		Priority:     6,
		Required:     false,
		Type:         dynamicfields.MultiLineTextType,
		DefaultValue: "",
		HelpText:     i18n.M(ctx, i18n.K.IntegrationDockerFieldsSwarmDnsResolversHelp),
		Conditions: []dynamicfields.Condition{
			{
				ParentField: SwarmModeFieldID,
				Value:       true,
			},
			{
				ParentField: SwarmServiceMeshFieldID,
				Value:       true,
			},
		},
	}

	useContainerNameAsID := dynamicfields.DynamicField{
		ID:           UseContainerNameAsIDFieldID,
		Description:  i18n.M(ctx, i18n.K.IntegrationDockerFieldsUseContainerNameAsId),
		Priority:     5,
		Required:     true,
		Type:         dynamicfields.BooleanType,
		DefaultValue: false,
		HelpText:     i18n.M(ctx, i18n.K.IntegrationDockerFieldsUseContainerNameAsIdHelp),
		Conditions: []dynamicfields.Condition{{
			ParentField: SwarmModeFieldID,
			Value:       false,
		}},
	}

	proxyURL := dynamicfields.DynamicField{
		ID:           ProxyURLFieldID,
		Description:  i18n.M(ctx, i18n.K.IntegrationDockerFieldsProxyUrl),
		Priority:     6,
		Required:     false,
		Type:         dynamicfields.URLType,
		DefaultValue: "",
		HelpText:     i18n.M(ctx, i18n.K.IntegrationDockerFieldsProxyUrlHelp),
		Conditions: []dynamicfields.Condition{{
			ParentField: SwarmModeFieldID,
			Value:       false,
		}},
	}

	return []dynamicfields.DynamicField{
		connectionMode,
		socketPath,
		hostURL,
		swarmMode,
		swarmServiceMesh,
		swarmDNSResolvers,
		useContainerNameAsID,
		proxyURL,
	}
}
