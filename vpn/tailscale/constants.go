package tailscale

import (
	"github.com/aws/smithy-go/ptr"

	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

const (
	authKeyFieldName        = "authKey"
	coordinatorUrlFieldName = "coordinatorUrl"
)

var configurationFields = []*dynamic_fields.DynamicField{
	{
		ID:          authKeyFieldName,
		Priority:    0,
		Description: "Tailscale auth key",
		Required:    true,
		Sensitive:   true,
		Type:        dynamic_fields.SingleLineTextType,
	},
	{
		ID:          coordinatorUrlFieldName,
		Priority:    1,
		Description: "Tailnet coordinator URL",
		Required:    false,
		Sensitive:   false,
		Type:        dynamic_fields.URLType,
		HelpText:    ptr.String("Custom coordinator server URL. Leave empty to use the default (tailscale.com)."),
	},
}
