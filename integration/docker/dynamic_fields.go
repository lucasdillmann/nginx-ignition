package docker

import "dillmann.com.br/nginx-ignition/core/common/dynamicfields"

var (
	connectionModeField = dynamicfields.DynamicField{
		ID:          "connectionMode",
		Description: "Connection mode",
		Priority:    1,
		Required:    true,
		Type:        dynamicfields.EnumType,
		EnumOptions: &[]*dynamicfields.EnumOption{
			{ID: "SOCKET", Description: "Socket"},
			{ID: "TCP", Description: "TCP"},
		},
		DefaultValue: stringPtr("SOCKET"),
	}

	socketPathField = dynamicfields.DynamicField{
		ID:           "socketPath",
		Description:  "Socket path",
		Priority:     2,
		Required:     true,
		Type:         dynamicfields.SingleLineTextType,
		DefaultValue: stringPtr("/var/run/docker.sock"),
		Condition: &dynamicfields.Condition{
			ParentField: connectionModeField.ID,
			Value:       "SOCKET",
		},
	}

	hostUrlField = dynamicfields.DynamicField{
		ID:          "hostUrl",
		Description: "Host URL",
		Priority:    2,
		Required:    true,
		Type:        dynamicfields.URLType,
		HelpText:    stringPtr("The URL to be used to connect to the Docker daemon, such as tcp://example.com:2375"),
		Condition: &dynamicfields.Condition{
			ParentField: connectionModeField.ID,
			Value:       "TCP",
		},
	}

	proxyUrlField = dynamicfields.DynamicField{
		ID:          "proxyUrl",
		Description: "Apps URL",
		Priority:    3,
		Required:    false,
		Type:        dynamicfields.URLType,
		HelpText: stringPtr("The URL to be used when proxying a request to a Docker container. Use this if the " +
			"apps are exposed in another address that isn't the container IP (from the Docker network)."),
	}
)

func stringPtr(s string) *string {
	return &s
}
