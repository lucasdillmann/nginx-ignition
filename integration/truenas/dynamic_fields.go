package truenas

import "dillmann.com.br/nginx-ignition/core/common/dynamic_fields"

var (
	url = dynamic_fields.DynamicField{
		ID:          "url",
		Description: "URL",
		Priority:    1,
		Required:    true,
		HelpText: stringPtr("The URL where your NAS is accessible, like http://192.168.0.2 or " +
			"https://nas.yourdomain.com"),
		Type: dynamic_fields.URLType,
	}

	proxyUrl = dynamic_fields.DynamicField{
		ID:          "proxyUrl",
		Description: "Apps URL",
		Priority:    2,
		Required:    false,
		HelpText: stringPtr("The URL to be used when proxying a request to a TrueNAS app. Use this if the apps " +
			"are exposed in another address that isn't the same as the main URL above."),
		Type: dynamic_fields.URLType,
	}

	username = dynamic_fields.DynamicField{
		ID:          "username",
		Description: "Username",
		Priority:    3,
		Required:    true,
		Sensitive:   false,
		Type:        dynamic_fields.SingleLineTextType,
	}

	password = dynamic_fields.DynamicField{
		ID:          "password",
		Description: "Password",
		Priority:    4,
		Required:    true,
		Sensitive:   true,
		Type:        dynamic_fields.SingleLineTextType,
	}
)

func stringPtr(s string) *string {
	return &s
}
