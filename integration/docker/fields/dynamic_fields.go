package fields

import (
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/ptr"
)

var (
	ConnectionMode = dynamicfields.DynamicField{
		ID:           "connectionMode",
		Description:  "Connection mode",
		Priority:     1,
		Required:     true,
		Type:         dynamicfields.EnumType,
		DefaultValue: SocketConnectionMode,
		EnumOptions: []dynamicfields.EnumOption{
			{ID: SocketConnectionMode, Description: "Socket"},
			{ID: TCPConnectionMode, Description: "TCP"},
		},
	}

	SocketPath = dynamicfields.DynamicField{
		ID:           "socketPath",
		Description:  "Socket path",
		Priority:     2,
		Required:     true,
		Type:         dynamicfields.SingleLineTextType,
		DefaultValue: "/var/run/docker.sock",
		Conditions: []dynamicfields.Condition{{
			ParentField: ConnectionMode.ID,
			Value:       SocketConnectionMode,
		}},
	}

	HostURL = dynamicfields.DynamicField{
		ID:           "hostUrl",
		Description:  "Host URL",
		Priority:     3,
		Required:     true,
		Type:         dynamicfields.URLType,
		DefaultValue: "",
		HelpText:     ptr.Of("The URL to be used to connect to Docker (such as tcp://example.com:2375)"),
		Conditions: []dynamicfields.Condition{{
			ParentField: ConnectionMode.ID,
			Value:       TCPConnectionMode,
		}},
	}

	SwarmMode = dynamicfields.DynamicField{
		ID:           "swarmMode",
		Description:  "Swarm mode",
		Priority:     4,
		Required:     true,
		Type:         dynamicfields.BooleanType,
		DefaultValue: false,
		HelpText: ptr.Of("When enabled, ignition will retrieve the available options by looking for the " +
			"deployed Swarm services instead of resolving available containers"),
	}

	SwarmServiceMesh = dynamicfields.DynamicField{
		ID:           "swarmServiceMesh",
		Description:  "Service mesh",
		Priority:     5,
		Required:     true,
		Type:         dynamicfields.BooleanType,
		DefaultValue: false,
		HelpText: ptr.Of("When enabled, nginx will be configured to reach Swarm services using the service mesh " +
			"(internal DNS names) when an ingress is selected as the proxy target"),
		Conditions: []dynamicfields.Condition{{
			ParentField: SwarmMode.ID,
			Value:       true,
		}},
	}

	SwarmDNSResolvers = dynamicfields.DynamicField{
		ID:           "swarmDnsResolvers",
		Description:  "Swarm DNS resolvers",
		Priority:     6,
		Required:     false,
		Type:         dynamicfields.MultiLineTextType,
		DefaultValue: "",
		HelpText: ptr.Of("Overrides the default DNS resolvers used by nginx when resolving Swarm services. " +
			"One IP address per line."),
		Conditions: []dynamicfields.Condition{
			{
				ParentField: SwarmMode.ID,
				Value:       true,
			},
			{
				ParentField: SwarmServiceMesh.ID,
				Value:       true,
			},
		},
	}

	UseContainerNameAsID = dynamicfields.DynamicField{
		ID:           "useContainerNameAsId",
		Description:  "Use container name as ID",
		Priority:     5,
		Required:     true,
		Type:         dynamicfields.BooleanType,
		DefaultValue: false,
		HelpText: ptr.Of("When enabled, ignition will use the container name as the ID instead of the " +
			"container's actual ID. Use this option when the containers are recreated constantly and/or managed " +
			"by a third-party tool."),
		Conditions: []dynamicfields.Condition{{
			ParentField: SwarmMode.ID,
			Value:       false,
		}},
	}

	ProxyURL = dynamicfields.DynamicField{
		ID:           "proxyUrl",
		Description:  "Proxy URL",
		Priority:     6,
		Required:     false,
		Type:         dynamicfields.URLType,
		DefaultValue: "",
		HelpText: ptr.Of("The URL to be used when proxying a request to a Docker container using a port " +
			"exposed on the host. If not set, the container IP will be used instead."),
		Conditions: []dynamicfields.Condition{{
			ParentField: SwarmMode.ID,
			Value:       false,
		}},
	}

	All = []dynamicfields.DynamicField{
		ConnectionMode,
		SocketPath,
		HostURL,
		SwarmMode,
		SwarmServiceMesh,
		SwarmDNSResolvers,
		UseContainerNameAsID,
		ProxyURL,
	}
)
