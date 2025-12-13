package tailscale

import (
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/ptr"
)

const (
	authKeyFieldName        = "authKey"
	coordinatorUrlFieldName = "coordinatorUrl"
)

var configurationFields = []*dynamicfields.DynamicField{
	{
		ID:          authKeyFieldName,
		Priority:    0,
		Description: "Tailscale auth key",
		Required:    true,
		Sensitive:   true,
		Type:        dynamicfields.SingleLineTextType,
	},
	{
		ID:          coordinatorUrlFieldName,
		Priority:    1,
		Description: "Tailnet coordinator URL",
		Required:    false,
		Sensitive:   false,
		Type:        dynamicfields.URLType,
		HelpText:    ptr.Of("Custom coordinator server URL. Leave empty to use the default (tailscale.com)."),
	},
}

var importantInstructions = []string{
	"An auth key can be generated in the Tailscale Admin console under Settings > Personal settings.",
	"When generating the key, make sure to generate a Reusable, Ephemeral and Pre-approved key. " +
		"Otherwise, nginx ignition will not be able to property manage and register virtual devices in the network.",
	"nginx ignition will use Tailscale to automatically provision server certificates for your ts.net domain. Make sure " +
		"that such possibility is enabled under Admin console > DNS > HTTP certificates.",
}
