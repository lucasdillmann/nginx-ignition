package docker

import "dillmann.com.br/nginx-ignition/core/common/dynamic_fields"

var (
	connectionModeField = dynamic_fields.DynamicField{
		ID:          "connectionMode",
		Description: "Connection mode",
		Priority:    1,
		Required:    true,
		Type:        dynamic_fields.EnumType,
		EnumOptions: &[]*dynamic_fields.EnumOption{
			{ID: "SOCKET", Description: "Socket"},
			{ID: "TCP", Description: "TCP"},
		},
		DefaultValue: stringPtr("SOCKET"),
	}

	socketPathField = dynamic_fields.DynamicField{
		ID:           "socketPath",
		Description:  "Socket path",
		Priority:     2,
		Required:     true,
		Type:         dynamic_fields.SingleLineTextType,
		DefaultValue: stringPtr("/var/run/docker.sock"),
		Condition: &dynamic_fields.Condition{
			ParentField: connectionModeField.ID,
			Value:       "SOCKET",
		},
	}

	hostUrlField = dynamic_fields.DynamicField{
		ID:          "hostUrl",
		Description: "Host URL",
		Priority:    2,
		Required:    true,
		Type:        dynamic_fields.URLType,
		HelpText:    stringPtr("The URL to be used to connect to the Docker daemon, such as tcp://example.com:2375"),
		Condition: &dynamic_fields.Condition{
			ParentField: connectionModeField.ID,
			Value:       "TCP",
		},
	}

	proxyUrlField = dynamic_fields.DynamicField{
		ID:          "proxyUrl",
		Description: "Apps URL",
		Priority:    3,
		Required:    false,
		Type:        dynamic_fields.URLType,
		HelpText: stringPtr("The URL to be used when proxying a request to a Docker container. Use this if the " +
			"apps are exposed in another address that isn't the container IP (from the Docker network)."),
	}
)

func stringPtr(s string) *string {
	return &s
}
