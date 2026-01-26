package tailscale

import (
	"context"

	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

const (
	authKeyFieldName        = "authKey"
	coordinatorURLFieldName = "coordinatorUrl"
)

func configurationFields(ctx context.Context) []dynamicfields.DynamicField {
	return []dynamicfields.DynamicField{
		{
			ID:          authKeyFieldName,
			Priority:    0,
			Description: i18n.M(ctx, i18n.K.VpnTailscaleAuthKey),
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          coordinatorURLFieldName,
			Priority:    1,
			Description: i18n.M(ctx, i18n.K.VpnTailscaleTailnetCoordinatorUrl),
			Required:    false,
			Sensitive:   false,
			Type:        dynamicfields.URLType,
			HelpText:    i18n.M(ctx, i18n.K.VpnTailscaleCoordinatorUrlHelp),
		},
	}
}

func importantInstructions(ctx context.Context) []*i18n.Message {
	return []*i18n.Message{
		i18n.M(ctx, i18n.K.VpnTailscaleInstructionAuthKey),
		i18n.M(ctx, i18n.K.VpnTailscaleInstructionKeySettings),
		i18n.M(ctx, i18n.K.VpnTailscaleInstructionSsl),
	}
}
