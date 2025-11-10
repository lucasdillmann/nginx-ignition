package truenas

import "dillmann.com.br/nginx-ignition/core/common/dynamicfields"

var (
	urlField = dynamicfields.DynamicField{
		ID:          "url",
		Description: "URL",
		Priority:    1,
		Required:    true,
		HelpText: stringPtr("The URL where your NAS is accessible, like http://192.168.0.2 or " +
			"https://nas.yourdomain.com"),
		Type: dynamicfields.URLType,
	}

	proxyUrlField = dynamicfields.DynamicField{
		ID:          "proxyUrl",
		Description: "Apps URL",
		Priority:    2,
		Required:    false,
		HelpText: stringPtr("The URL to be used when proxying a request to a TrueNAS app. Use this if the apps " +
			"are exposed in another address that isn't the same as the main URL above."),
		Type: dynamicfields.URLType,
	}

	usernameField = dynamicfields.DynamicField{
		ID:          "username",
		Description: "Username",
		Priority:    3,
		Required:    true,
		Sensitive:   false,
		Type:        dynamicfields.SingleLineTextType,
	}

	passwordField = dynamicfields.DynamicField{
		ID:          "password",
		Description: "Password",
		Priority:    4,
		Required:    true,
		Sensitive:   true,
		Type:        dynamicfields.SingleLineTextType,
	}
)

func stringPtr(s string) *string {
	return &s
}
