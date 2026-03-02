package fields

import (
	"context"

	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

const (
	URLFieldID       = "url"
	ProxyURLFieldID  = "proxyUrl"
	UsernameFieldID  = "username"
	PasswordFieldID  = "password"
	LegacyAPIFieldID = "legacyApi"
)

func DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return []dynamicfields.DynamicField{
		{
			ID:          URLFieldID,
			Description: i18n.M(ctx, i18n.K.IntegrationTruenasUrl),
			Priority:    1,
			Required:    true,
			HelpText:    i18n.M(ctx, i18n.K.IntegrationTruenasUrlHelp),
			Type:        dynamicfields.URLType,
		},
		{
			ID:           LegacyAPIFieldID,
			Description:  i18n.M(ctx, i18n.K.IntegrationTruenasLegacyApi),
			HelpText:     i18n.M(ctx, i18n.K.IntegrationTruenasLegacyApiHelp),
			Priority:     2,
			Required:     true,
			Type:         dynamicfields.BooleanType,
			DefaultValue: false,
		},
		{
			ID:          ProxyURLFieldID,
			Description: i18n.M(ctx, i18n.K.IntegrationTruenasProxyUrl),
			Priority:    3,
			Required:    false,
			HelpText:    i18n.M(ctx, i18n.K.IntegrationTruenasProxyUrlHelp),
			Type:        dynamicfields.URLType,
		},
		{
			ID:          UsernameFieldID,
			Description: i18n.M(ctx, i18n.K.CommonUsername),
			Priority:    4,
			Required:    true,
			Sensitive:   false,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          PasswordFieldID,
			Description: i18n.M(ctx, i18n.K.CommonPassword),
			Priority:    5,
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
	}
}
