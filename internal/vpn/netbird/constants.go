package netbird

import (
	"context"

	"github.com/lucasdillmann/nginx-ignition/internal/core/common/dynamicfields"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/i18n"
)

const (
	setupKeyFieldName      = "setupKey"
	managementURLFieldName = "managementUrl"
)

func configurationFields(ctx context.Context) []dynamicfields.DynamicField {
	return []dynamicfields.DynamicField{
		{
			ID:          setupKeyFieldName,
			Priority:    0,
			Description: i18n.M(ctx, i18n.K.VpnNetbirdSetupKey),
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          managementURLFieldName,
			Priority:    1,
			Description: i18n.M(ctx, i18n.K.VpnNetbirdManagementUrl),
			Required:    false,
			Sensitive:   false,
			Type:        dynamicfields.URLType,
			HelpText:    i18n.M(ctx, i18n.K.VpnNetbirdManagementUrlHelp),
		},
	}
}

func importantInstructions(ctx context.Context) []*i18n.Message {
	return []*i18n.Message{
		i18n.M(ctx, i18n.K.VpnNetbirdInstructionSetupKey),
		i18n.M(ctx, i18n.K.VpnNetbirdInstructionKeySettings),
	}
}
