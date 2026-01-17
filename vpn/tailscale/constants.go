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
			Description: i18n.M(ctx, i18n.K.TailscaleCommonAuthKey),
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          coordinatorURLFieldName,
			Priority:    1,
			Description: i18n.M(ctx, i18n.K.TailscaleCommonTailnetCoordinatorUrl),
			Required:    false,
			Sensitive:   false,
			Type:        dynamicfields.URLType,
			HelpText:    i18n.M(ctx, i18n.K.TailscaleCommonCoordinatorUrlHelp),
		},
	}
}

func importantInstructions(ctx context.Context) []*i18n.Message {
	return []*i18n.Message{
		i18n.M(ctx, i18n.K.TailscaleCommonInstructionAuthKey),
		i18n.M(ctx, i18n.K.TailscaleCommonInstructionKeySettings),
		i18n.M(ctx, i18n.K.TailscaleCommonInstructionSsl),
	}
}
