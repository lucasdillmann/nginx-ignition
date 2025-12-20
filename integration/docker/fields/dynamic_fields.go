package fields

import (
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/ptr"
)

var (
	ConnectionMode = dynamicfields.DynamicField{
		ID:          "connectionMode",
		Description: "Connection mode",
		Priority:    1,
		Required:    true,
		Type:        dynamicfields.EnumType,
		EnumOptions: &[]*dynamicfields.EnumOption{
			{ID: "SOCKET", Description: "Socket"},
			{ID: "TCP", Description: "TCP"},
		},
		DefaultValue: ptr.Of("SOCKET"),
	}

	SocketPath = dynamicfields.DynamicField{
		ID:           "socketPath",
		Description:  "Socket path",
		Priority:     2,
		Required:     true,
		Type:         dynamicfields.SingleLineTextType,
		DefaultValue: ptr.Of("/var/run/docker.sock"),
		Condition: &dynamicfields.Condition{
			ParentField: ConnectionMode.ID,
			Value:       "SOCKET",
		},
	}

	HostURL = dynamicfields.DynamicField{
		ID:          "hostUrl",
		Description: "Host URL",
		Priority:    3,
		Required:    true,
		Type:        dynamicfields.URLType,
		HelpText:    ptr.Of("The URL to be used to connect to the Docker daemon, such as tcp://example.com:2375"),
		Condition: &dynamicfields.Condition{
			ParentField: ConnectionMode.ID,
			Value:       "TCP",
		},
	}

	SwarmMode = dynamicfields.DynamicField{
		ID:          "swarmMode",
		Description: "Swarm mode",
		Priority:    4,
		Required:    true,
		Type:        dynamicfields.BooleanType,
		HelpText: ptr.Of("When enabled, ignition will retrieve the available options by looking for the " +
			"deployed Swarm services instead of resolving available containers"),
	}

	SwarmServiceMesh = dynamicfields.DynamicField{
		ID:          "swarmServiceMesh",
		Description: "Service mesh",
		Priority:    5,
		Required:    false,
		Type:        dynamicfields.BooleanType,
		HelpText: ptr.Of("When enabled, nginx will be configured to reach Swarm services using the service mesh " +
			"(internal DNS names)."),
		Condition: &dynamicfields.Condition{
			ParentField: SwarmMode.ID,
			Value:       true,
		},
	}

	SwarmDNSResolvers = dynamicfields.DynamicField{
		ID:          "swarmDnsResolvers",
		Description: "Swarm DNS resolvers",
		Priority:    6,
		Required:    false,
		Type:        dynamicfields.MultiLineTextType,
		HelpText: ptr.Of("Overrides the default DNS resolvers used by nginx when resolving Swarm services (" +
			"if omitted, nginx will use the default resolvers). Inform one resolver IP address per line."),
		Condition: &dynamicfields.Condition{
			ParentField: SwarmMode.ID,
			Value:       true,
		},
	}

	ProxyURL = dynamicfields.DynamicField{
		ID:          "proxyUrl",
		Description: "Proxy URL",
		Priority:    6,
		Required:    false,
		Type:        dynamicfields.URLType,
		HelpText: ptr.Of("The URL to be used when proxying a request to a Docker container using a port " +
			"exposed on the host. If not set, the container IP will be used instead."),
		Condition: &dynamicfields.Condition{
			ParentField: SwarmMode.ID,
			Value:       false,
		},
	}

	All = []*dynamicfields.DynamicField{
		&ConnectionMode,
		&SocketPath,
		&HostURL,
		&SwarmMode,
		&SwarmServiceMesh,
		&SwarmDNSResolvers,
		&ProxyURL,
	}
)
